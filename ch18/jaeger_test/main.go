package main

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.189.128:6831",
		},
		ServiceName: "shop",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))

	if err != nil {
		panic(err)
	}
	defer closer.Close()

	parentSpan := tracer.StartSpan("main")

	span := tracer.StartSpan("funcA", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 500)
	span.Finish()

	span2 := tracer.StartSpan("funcB", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 100)
	span2.Finish()

	parentSpan.Finish()
}
