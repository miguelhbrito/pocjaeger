package main

import (
	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/midtracing"
	"github.com/pocjaeger/tracing"
	stdlog "log"
	"net/http"
	"os"
)

func myTracingHandler(w http.ResponseWriter, r *http.Request) {
	tracer, closer := tracing.InitJaeger("serverTracing")
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	defer closer.Close()
	span := tracer.StartSpan("server-tracing", ext.RPCServerOption(spanCtx))
	defer span.Finish()
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/tracing", myTracingHandler)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)

	loggingMiddleware := midtracing.LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(router)

	if err := http.ListenAndServe(":8000", loggedRouter); err != nil {
		os.Exit(1)
	}
}
