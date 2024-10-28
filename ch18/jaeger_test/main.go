package main

import (
	"time"

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
			LocalAgentHostPort: "172.16.89.133:6831",
		},
		ServiceName: "alishop",
	}
	cfg.InitGlobalTracer("alishop")
	cfg.InitGlobalTracer("go-grpc-web")

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	defer closer.Close()

	span := tracer.StartSpan("go-grpc-web")
	time.Sleep(time.Second)
	defer span.Finish()
}
