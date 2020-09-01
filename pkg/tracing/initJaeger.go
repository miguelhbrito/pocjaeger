package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/zipkin"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func InitJaeger(service string) io.Closer {
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	cfg := config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			//already in default
			SamplingServerURL: "http://localhost:5778/sampling",
			Type:              "const",
			Param:             1.0,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Millisecond,
			LocalAgentHostPort:  "127.0.0.1:6831",
		},
	}

	jLogger := jaegerlog.StdLogger

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jLogger),
		config.ZipkinSharedRPCSpan(true),
		config.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		config.Extractor(opentracing.HTTPHeaders, zipkinPropagator))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}

func NewTracedRequest(method string, url string, span opentracing.Span) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, method)
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	return req, err
}
