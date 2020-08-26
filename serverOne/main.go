package main

import (
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pocjaeger/pkg/serverOne"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"net/http"
)

func TracingMiddleware (operationName string) func (next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			span := opentracing.GlobalTracer().StartSpan(operationName)
			r.WithContext(opentracing.ContextWithSpan(r.Context(), span))

			defer span.Finish()

			tid := tracing.GetTraceID(span)

			span.LogFields(
				otlog.String("event", "event-from-client"),
				otlog.String("TraceID", tid),
			)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func main() {
	closer := tracing.InitJaeger("server-one-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	middleware := TracingMiddleware("Request on Server One")
	http.Handle("/", middleware(http.HandlerFunc(serverOne.MyTracingHandlerServerOne)))
	http.Handle("/serverTwoResponse", http.HandlerFunc(serverOne.HandlerServerTwoResponse))

	log.Info().Msg("Server one listening on 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Error().Err(err).Msg("handle error")
	}
}