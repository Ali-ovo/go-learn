package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-learn/ch19/log/log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

const (
	traceName = "shop-otel"
)

var tp *trace.TracerProvider

func tracerProvider() error {
	// 创建 Jaeger 导出器
	url := "http://192.168.189.128:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp = trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL, // 必填参数 可以指定是多少版本
			semconv.ServiceNameKey.String("shop-user"),
			attribute.String("environment", "dev"),
			attribute.Int("ID", 1),
		)),
	)
	// 设置全局的提取器
	otel.SetTracerProvider(tp)
	// 设置全局的 传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func funcA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := otel.Tracer(traceName)
	spanCtx, span := tr.Start(ctx, "func-a")
	span.SetAttributes(attribute.String("name", "funcA"))
	fmt.Println("trace:", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())

	type _LogStruct struct {
		CurrentTime time.Time `json:"current_time"`
		PassWho     string    `json:"pass_who"`
		Name        string    `json:"name"`
	}

	log.InfofC(spanCtx, "this is funca log")

	logTest := _LogStruct{
		CurrentTime: time.Now(),
		PassWho:     "czc1",
		Name:        "func-a",
	}
	b, _ := json.Marshal(logTest)
	// 日志会发送到 jaeger
	log.InfofC(spanCtx, "this is funca log: %s", string(b))

	span.SetAttributes(attribute.Key("这是测试日志的key").String(string(b)))
	time.Sleep(time.Second)
	span.End()
}

func funcB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := otel.Tracer(traceName)
	spanCtx, span := tr.Start(ctx, "func-b")
	fmt.Println("trace:", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	time.Sleep(time.Second)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8090/server")
	req.Header.SetMethod("GET")

	// 拿到 p 传播器
	p := otel.GetTextMapPropagator()

	// 包裹
	headers := make(map[string]string)
	// 从 spanCtx 中 解析出相关信息到 propagation.MapCarrier(headers) 中
	// propagation.MapCarrier 是相关 spanCtx 信息到 headers 中   propagation.HeaderCarrier 注入到 http.Header 中
	//client := &http.Client{}
	//hreq, _ := http.NewRequest("GET", "http://127.0.0.1:8090/server", nil)
	//p.Inject(spanCtx, propagation.HeaderCarrier(hreq.Header))
	p.Inject(spanCtx, propagation.MapCarrier(headers)) // traceparent -> 00-f0e13ee095f28266f1660b1a9681c8c3-33ed0a659cf0b685-01

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fclient := fasthttp.Client{}
	fres := fasthttp.Response{}
	_ = fclient.Do(req, &fres)
	log.InfofC(spanCtx, "this is funcb log: %s", "imooc")

	span.End()
}

func main() {
	_ = tracerProvider()
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
	time.Sleep(time.Second * time.Duration(2))
	wg.Wait()

	span.End()
}
