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
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	tracesdk "go.opentelemetry.io/otel/trace"
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
	// 设置全局的传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func Server(c *gin.Context) {
	// 负载 span 的抽取和生成
	ctx := c.Request.Context()
	p := otel.GetTextMapPropagator()
	tr := tp.Tracer(traceName)
	// 从 propagation.HeaderCarrier(c.Request.Header) 中 解析出相关信息到 ctx
	sctx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
	_, span := tr.Start(tracesdk.ContextWithRemoteSpanContext(sctx, tracesdk.SpanContextFromContext(sctx)), "gin-server")
	time.Sleep(time.Duration(500) * time.Millisecond)
	span.End()
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func main() {
	_ = tracerProvider()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {})
	r.GET("/server", Server)
	r.Run(":8090")
}
