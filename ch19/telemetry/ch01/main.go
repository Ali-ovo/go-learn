package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func main() {
	url := "http://172.16.89.133:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
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

	ctx, cancel := context.WithCancel(context.Background())

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)

	tr := otel.Tracer("shop-otel")
	_, span := tr.Start(ctx, "func-main")

	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.String("key1", "value1"))
	attrs = append(attrs, attribute.Bool("key2", true))
	attrs = append(attrs, attribute.Int("key3", 100))

	span.SetAttributes(attrs...)
	time.Sleep(1 * time.Second)

	span.End()

}
