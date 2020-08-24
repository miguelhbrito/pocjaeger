package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/uber/jaeger-client-go"
	"log"
	"net/http"
	"os"
)

func myTracingHandlerServerTwo(w http.ResponseWriter, r *http.Request) {
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.GlobalTracer().StartSpan("server-two-operation", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	tid := span.Context().(jaeger.SpanContext).TraceID()
	tids := fmt.Sprintf("%s", tid)

	log.Println("TraceID to see the tracing of this request: ", tids)

	span.LogFields(
		otlog.String("event", "event-from-client"),
		otlog.String("traceID", tids),
	)

	w.Header().Set("traceID", tids)
}

func main() {
	closer := tracing.InitJaeger("server-two-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/server-two", http.HandlerFunc(myTracingHandlerServerTwo))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(1)
	}
}
