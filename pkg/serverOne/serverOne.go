package serverOne

import (
	"github.com/opentracing/opentracing-go"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"net/http"
)

func MyTracingHandlerServerOne(w http.ResponseWriter, r *http.Request) {
	//rootSpan := opentracing.GlobalTracer().StartSpan("Request on Server One")
	//ctx := opentracing.ContextWithSpan(context.Background(), rootSpan)
	//defer rootSpan.Finish()

	ctx := r.Context()
	response, err := client.DoRequest(ctx)
	if err != nil {
		log.Error().Err(err).Msg("handle error")
		return
	}

	//body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("handle error")
		return
	}
	defer response.Body.Close()

	tid := tracing.GetTraceID(opentracing.SpanFromContext(ctx))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("traceID", tid)
	_, err = w.Write([]byte("faslkfhskdjhfskjdhfasf"))
	if err != nil {
		log.Error().Err(err).Msg("handle error")
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
			log.Error().Err(err).Msg("handle error")
		}
	}
}
