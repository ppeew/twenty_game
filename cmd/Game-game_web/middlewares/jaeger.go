package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func Jaeger() gin.HandlerFunc {
	return func(context *gin.Context) {
		cfg := jaegercfg.Configuration{
			ServiceName: "game_web",
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: "139.159.234.134:6831",
			},
		}

		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
		if err != nil {
			panic(err)
		}
		defer closer.Close()

		parentSpan := tracer.StartSpan(context.Request.URL.Path)
		defer parentSpan.Finish()

		context.Set("tracer", tracer)
		context.Set("parentSpan", parentSpan)
		context.Next()
	}
}
