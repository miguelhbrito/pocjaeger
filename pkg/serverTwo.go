package main

import (
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func myTracingHandlerServerTwo(w http.ResponseWriter, r *http.Request) {
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.GlobalTracer().StartSpan("server-two-operation", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	tid := tracing.GetTraceID(span)
	sid := tracing.GetSpanID(span)

	span.SetTag("traceID", tid)
	span.LogKV("event", "request from server one")
	span.LogKV("request status", http.StatusOK)

	w.WriteHeader(200)
	w.Header().Set("traceID", tid)
	w.Header().Set("spanID", sid)

	err := json.NewEncoder(w).Encode(client.Response{
		TraceID: tid,
		SpanID:  sid,
	})
	if err != nil {
		log.Error().Err(err)
	}
}

func main() {
	closer := tracing.InitJaeger("server-two-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/server-two", http.HandlerFunc(myTracingHandlerServerTwo))

	log.Info().Msg("Server two listening on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(1)
	}
}
