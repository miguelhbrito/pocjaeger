package midtracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		closer := tracing.InitJaeger("server-one-tracing")
		defer closer.Close()
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := opentracing.GlobalTracer().StartSpan("server-one-midtracing", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		tid := span.Context().(jaeger.SpanContext).TraceID()
		tids := fmt.Sprintf("%s", tid)

		span.LogFields(
			otlog.String("event", "event-from-client"),
			otlog.String("TraceID", tids),
		)
		w.Write([]byte(tids))

		next.ServeHTTP(w, r)
	})
}
