package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	traceName = "shop-otel"
)

var tp *traceSdk.TracerProvider

func tracerProvider(url string) error {
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

	if err != nil {
		panic(err)
	}

	tp = traceSdk.NewTracerProvider(
		traceSdk.WithBatcher(jexp),
		traceSdk.WithResource(
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

// func Server(c *gin.Context) {
// 	// span的抽取和生成
// 	ctx := c.Request.Context()
// 	p := otel.GetTextMapPropagator()

// 	tr := tp.Tracer(traceName)
// 	sctx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
// 	_, span := tr.Start(sctx, "server")
// 	time.Sleep(1 * time.Second)
// 	span.End()
// 	c.JSON(200, gin.H{})
// }

func Server(c *gin.Context) {
	// span的抽取和生成
	ctx := c.Request.Context()
	// p := otel.GetTextMapPropagator()
	// tr := tp.Tracer(traceName)
	// sctx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

	tr := tp.Tracer(traceName)
	traceID := c.Request.Header.Get("trace-id")
	spanID := c.Request.Header.Get("span-id")

	traceid, _ := trace.TraceIDFromHex(traceID)
	spanid, _ := trace.SpanIDFromHex(spanID)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID(traceid),
		SpanID:     trace.SpanID(spanid),
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})

	carrier := propagation.HeaderCarrier{}
	carrier.Set("trace-id", traceID)
	carrier.Set("span-id", spanID)
	propagator := otel.GetTextMapPropagator()
	pctx := propagator.Extract(ctx, carrier)
	sctx := trace.ContextWithRemoteSpanContext(pctx, spanCtx)

	_, span := tr.Start(sctx, "server")
	time.Sleep(1 * time.Second)
	span.End()
	c.JSON(200, gin.H{})
}

func main() {
	url := "http://172.16.89.133:14268/api/traces"
	tracerProvider(url)

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {

	})

	r.GET("/server", Server)

	r.Run(":8090")
}
