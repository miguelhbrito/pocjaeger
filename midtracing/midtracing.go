package midtracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/tracing"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer, closer := tracing.InitJaeger("server-tracing")
		defer closer.Close()
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("server-midtracing", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			span.SetTag("spanID", sc.SpanID())
			span.SetTag("traceID", sc.TraceID())
		}

		next.ServeHTTP(w, r)
	})
}
