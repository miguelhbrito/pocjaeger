package serverTwo

import (
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"net/http"
)

func MyTracingHandlerServerTwo(w http.ResponseWriter, r *http.Request) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.GlobalTracer().StartSpan("Server two request handler", opentracing.ChildOf(spanCtx))
	defer span.Finish()

	tid := tracing.GetTraceID(span)

	span.SetTag("traceID", tid)
	span.SetTag("request method", r.Method)
	span.SetTag("Server", "TWO")

	span.LogKV("event", "request made by server one")
	span.LogKV("status code", http.StatusOK)

	w.WriteHeader(200)
	w.Header().Set("traceID", tid)

	spanResponse := opentracing.GlobalTracer().StartSpan("Writing response to server one", opentracing.ChildOf(span.Context()))
	spanResponse.SetTag("status code", http.StatusOK)
	defer spanResponse.Finish()
	err = json.NewEncoder(w).Encode(client.Response{
		TraceID: tid,
	})
	if err != nil {
		log.Error().Err(err)
		ext.LogError(span, err)
		return
	}
}
