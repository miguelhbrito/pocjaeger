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
	rootSpan := opentracing.GlobalTracer().StartSpan("Request from Server One")
	log.Debug().Msgf("Trace ID: %s", tracing.GetSpanID(rootSpan))
	ctx := opentracing.ContextWithSpan(context.Background(), rootSpan)
	defer rootSpan.Finish()

	//returns := tracing.RunTracedFunction(client.DoRequest, ctx)
	//
	//response := returns[0].Interface().(*http.Response)
	//var err error
	//if returns[1].Interface() != nil {
	//	err = returns[1].Interface().(error)
	//}

	doRequestTraced := tracing.MakeTracedFunction(client.DoRequest).(func(context.Context) (*http.Response, error))
	response, err := doRequestTraced(ctx)

	if err != nil {
		ext.LogError(rootSpan, err)
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err)
		return
	}
	defer response.Body.Close()

	tid := tracing.GetTraceID(rootSpan)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("traceID", tid)
	_, err = w.Write(body)
	if err != nil {
		log.Error().Err(err)
		return
	}
}

func HandlerServerTwoResponse(w http.ResponseWriter, r *http.Request) {
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.GlobalTracer().StartSpan("Request from server Two", opentracing.ChildOf(spanCtx))
	defer span.Finish()
	span.SetTag("traceID", tracing.GetTraceID(span))
	span.SetTag("Request method", r.Method)
	span.SetTag("Server", "ONE")

	if r.Method != "GET" {
		w.WriteHeader(404)
	} else {
		span.LogKV("event", "Receiving request from server two")
		_, err := w.Write([]byte("Success on responding to server one"))
		if err != nil {
			log.Error().Err(err)
		}
	}
}
