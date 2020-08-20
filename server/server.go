package main

import (
	"github.com/pocjaeger/midtracing"
	"net/http"
	"os"
)

func myTracingHandler(w http.ResponseWriter, r *http.Request) {
	//tracer, closer := tracing.InitJaeger("serverTracing")
	//spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	//defer closer.Close()
	//span := tracer.StartSpan("server-tracing", ext.RPCServerOption(spanCtx))
	//jaeger.SpanID()
	//defer span.Finish()
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", myTracingHandler)

	tracedRouter := midtracing.TracingMiddleware(router)

	if err := http.ListenAndServe(":8000", tracedRouter); err != nil {
		os.Exit(1)
	}
}
