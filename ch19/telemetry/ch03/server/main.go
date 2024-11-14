package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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

func Server(c *gin.Context) {
	// span的抽取和生成
	ctx := c.Request.Context()
	p := otel.GetTextMapPropagator()

	tr := tp.Tracer(traceName)
	sctx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
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
