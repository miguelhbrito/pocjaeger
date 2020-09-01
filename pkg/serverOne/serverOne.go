package serverOne

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

func MyTracingHandlerServerOne(w http.ResponseWriter, r *http.Request) {
	rootSpan := opentracing.GlobalTracer().StartSpan("Server one request handler")
	log.Debug().Msgf("Trace ID: %s", tracing.GetSpanID(rootSpan))
	ctx := opentracing.ContextWithSpan(context.Background(), rootSpan)
	defer rootSpan.Finish()

	rootSpan.SetTag("traceID", tracing.GetTraceID(rootSpan))
	rootSpan.SetTag("request method", r.Method)
	rootSpan.SetTag("Server", "ONE")

	response, err := client.DoRequest(ctx)
	if err != nil {
		log.Error().Err(err)
		ext.LogError(rootSpan, err)
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err)
		ext.LogError(rootSpan, err)
		return
	}
	defer response.Body.Close()

	rootSpan.LogKV("event", "request received on server one")
	rootSpan.LogKV("status code", http.StatusOK)

	tid := tracing.GetTraceID(rootSpan)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("traceID", tid)
	w.WriteHeader(200)
	_, err = w.Write(body)
	if err != nil {
		log.Error().Err(err)
		ext.LogError(rootSpan, err)
		return
	}
}
