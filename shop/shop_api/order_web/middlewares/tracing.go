package middlewares

import (
	"fmt"
	"go-learn/shop/shop_api/order_web/global"

	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}

		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}

		defer closer.Close()

		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}
}
