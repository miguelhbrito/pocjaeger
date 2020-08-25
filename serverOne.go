package main

import (
	"context"
	"encoding/json"
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

	body, err := client.DoRequest(ctx)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}

	var tid string
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		tid = sc.TraceID().String()
	}

	log.Println("TraceID to see the tracing of this request: ", tid)

	bodyString, _ := json.Marshal(body)

	w.Header().Set("traceID", tid)
	w.Write(bodyString)
}

func main() {
	closer := tracing.InitJaeger("server-one-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/", http.HandlerFunc(myTracingHandlerServerOne))

	log.Println("Server one listening on 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		os.Exit(1)
	}
}
