package util

import (
	"io"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func MakeTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  os.Getenv("JAEGER_AGENT_HOST_PORT"),
		},
	}
	return cfg.NewTracer()
}
