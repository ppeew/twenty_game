package tests

import (
	"github.com/opentracing/opentracing-go"
	"time"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"testing"
)

func TestJaeger(t *testing.T) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "139.159.234.134:6831",
		},
		ServiceName: "twenty_game",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	parentSpan := tracer.StartSpan("game_web parent")
	span := tracer.StartSpan("game1", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Second)

	span2 := tracer.StartSpan("game2", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Second * 2)

	span.Finish()
	span2.Finish()
	parentSpan.Finish()
}
