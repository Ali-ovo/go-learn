package trace

import (
	"shop/gmicro/pkg/log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

/*
初始化不同的 export 的设置
*/

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
)

var (
	// set struct 空结构体 不占内存 zerobase
	agents = make(map[string]struct{})
	lock   sync.Mutex
)

func InitAgent(o Options) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := agents[o.Endpoint]
	if ok {
		return
	}
	err := startAgent(o)
	if err != nil {
		return
	}
	agents[o.Endpoint] = struct{}{}
}

func startAgent(o Options) error {
	var sexp sdktrace.SpanExporter
	var err error

	opts := []sdktrace.TracerProviderOption{
		// WithSampler() 设置采样策略
		// ParentBased() 根据上一个 Span 是否抽样来决定当前 Span 是否抽样
		// 使用采样什么策略
		// AlwaysSample(): 总是对 span 进行采样
		// NeverSample(): 总是对 span 不进行采样
		// TraceIdRatioBased(double): 基于 trace ID 的一定比例进行采样
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(o.Sampler))),
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(o.Name))), // 设置链路追踪 提供者 的名字
	}

	if len(o.Endpoint) > 0 {
		switch o.Batcher {
		case kindJaeger:
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
			if err != nil {
				return err
			}
		case kindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err
			}
		}
		opts = append(opts, sdktrace.WithBatcher(sexp))
	}

	tp := sdktrace.NewTracerProvider(opts...)
	// 设置全局提取器
	otel.SetTracerProvider(tp)
	// 设置全局的传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))
	return nil
}
