package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	traceName = "shop-otel"
)

var tp *trace.TracerProvider

func tracerProvider(url string) error {
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

	if err != nil {
		panic(err)
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("shop-user"),
				attribute.String("env", "dev"),
				attribute.Int("ID", 1),
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func funcA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer(traceName)
	_, span := tr.Start(ctx, "func-A")
	span.SetAttributes(attribute.String("func", "func-A"))

	type _LogStruct struct {
		CurrentTime time.Time `json:"current_time"`
		PassWho     string    `json:"pass_who"`
		Name        string    `json:"name"`
	}

	logTest := _LogStruct{
		CurrentTime: time.Now(),
		PassWho:     "Ali",
		Name:        "func-A",
	}

	b, _ := json.Marshal(logTest)
	span.SetAttributes(attribute.Key("测试日志的 key").String(string(b)))
	time.Sleep(time.Second)
	span.End()

}

func funcB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer(traceName)
	_, span := tr.Start(ctx, "func-B")
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	time.Sleep(time.Second)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost:8090/server")
	req.Header.SetMethod("GET")

	// p := otel.GetTextMapPropagator()
	// headers := make(map[string]string)

	// p.Inject(spanCtx, propagation.MapCarrier(headers))
	// for key, value := range headers {
	// 	req.Header.Set(key, value)
	// }

	req.Header.Set("trace-id", span.SpanContext().TraceID().String())
	req.Header.Set("span-id", span.SpanContext().SpanID().String())

	fclient := fasthttp.Client{}
	fres := fasthttp.Response{}
	_ = fclient.Do(req, &fres)

	span.End()
}

func main() {
	url := "http://192.168.189.128:14268/api/traces"

	tracerProvider(url)
	ctx, cancel := context.WithCancel(context.Background())

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	tr := otel.Tracer(traceName)
	spanCtx, span := tr.Start(ctx, "func-main")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go funcA(spanCtx, wg)
	go funcB(spanCtx, wg)

	span.AddEvent("this is an event")
	time.Sleep(1 * time.Second)
	wg.Wait()

	span.End()

}
