package main

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
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

func main() {
	_ = tracerProvider()
	cli := redis.NewClient(&redis.Options{
		Addr: "192.168.189.128:6379",
	})

	if err := redisotel.InstrumentTracing(cli); err != nil {
		panic(err)
	}

	// 以下都是声明 tracer 的方法
	//tp.Tracer("traceName")
	//otel.GetTracerProvider().Tracer("traceName")
	tr := otel.Tracer("traceName")

	spanCtx, span := tr.Start(context.Background(), "redis")
	cli.Set(spanCtx, "name", "czc", 0)
	span.End()
	tp.Shutdown(context.Background())
}
