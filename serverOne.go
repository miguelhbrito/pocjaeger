package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/uber/jaeger-client-go"
	"log"
	"net/http"
	"os"
)

func myTracingHandlerServerOne(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("request-to-server-two")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	defer span.Finish()

	_, err := client.DoRequest(ctx)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}

	tid := span.Context().(jaeger.SpanContext).TraceID()
	tids := fmt.Sprintf("%s", tid)

	log.Println("TraceID to see the tracing of this request: ", tids)

	w.Header().Set("traceID", tids)
	w.Write([]byte(tids))
}

func main() {
	closer := tracing.InitJaeger("server-one-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/", http.HandlerFunc(myTracingHandlerServerOne))

	if err := http.ListenAndServe(":8000", nil); err != nil {
		os.Exit(1)
	}
}
