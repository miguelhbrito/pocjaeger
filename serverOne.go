package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
)

func myTracingHandlerServerOne(w http.ResponseWriter, r *http.Request) {
	rootSpan := opentracing.GlobalTracer().StartSpan("Request to Server One")
	ctx := opentracing.ContextWithSpan(context.Background(), rootSpan)
	defer rootSpan.Finish()

	response, err := client.DoRequest(ctx)
	if err != nil {
		ext.LogError(rootSpan, err)
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err)
	}
	defer response.Body.Close()

	tid := tracing.GetTraceID(rootSpan)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("traceID", tid)
	_, err = w.Write(body)
	if err != nil {
		log.Error().Err(err)
	}
}

func main() {
	closer := tracing.InitJaeger("server-one-tracing")
	defer closer.Close()

	//TODO maybe change this to use midtracing
	http.Handle("/", http.HandlerFunc(myTracingHandlerServerOne))

	log.Info().Msg("Server one listening on 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		os.Exit(1)
	}
}
