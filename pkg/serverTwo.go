package main

import (
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/uber/jaeger-client-go"
	"log"
	"net/http"
	"os"
)

func myTracingHandlerServerTwo(w http.ResponseWriter, r *http.Request){
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.GlobalTracer().StartSpan("server-two-operation", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	var tid, sid string
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		tid = sc.TraceID().String()
		sid = sc.SpanID().String()
	}

	log.Println("TraceID to see the tracing of this request: ", tid)

	span.LogFields(
		otlog.String("event", "event-from-client"),
		otlog.String("traceID", tid),
		otlog.Int("request status", http.StatusOK),
	)

	w.WriteHeader(200)

	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(w.Header()),
	)

	w.Header().Set("traceID", tid)
	w.Header().Set("spanID", sid)

	_ = json.NewEncoder(w).Encode(struct {
		TraceID string `json:"trace-id"`
		SpanID  string `json:"span-id"`
	}{
		TraceID: tid,
		SpanID:  sid,
	})
}

func main() {
	closer := tracing.InitJaeger("server-two-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/server-two", http.HandlerFunc(myTracingHandlerServerTwo))

	log.Println("Server two listening on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(1)
	}
}
