package serverTwo

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/pkg/client"
	"github.com/pocjaeger/pkg/tracing"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

func DoRequest(ctx context.Context, url string) (*http.Response, error) {
	resp := &http.Response{}
	recover()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return resp, err
	}

	span, _ := opentracing.StartSpanFromContext(ctx, "Doing request")
	defer span.Finish()
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	err = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		return resp, err
	}

	resp, err = http.DefaultClient.Do(req)
	return resp, err
}

func MyTracingHandlerServerTwo(w http.ResponseWriter, r *http.Request) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Error().Err(err).Msg("handle error")
		return
	}
	span := opentracing.GlobalTracer().StartSpan("Request received on server two", opentracing.ChildOf(spanCtx))
	defer span.Finish()

	tid := tracing.GetTraceID(span)

	span.SetTag("traceID", tid)
	span.SetTag("Request method", r.Method)
	span.SetTag("Server", "TWO")

	span.LogKV("event", "request from server one")
	span.LogKV("request status", http.StatusOK)

	w.WriteHeader(200)
	w.Header().Set("traceID", tid)

	err = json.NewEncoder(w).Encode(client.Response{
		TraceID: tid,
	})
	if err != nil {
		log.Error().Err(err).Msg("handle error")
		return
	}

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	response, err := DoRequest(ctx, "http://localhost:8000/serverTwoResponse")
	if err != nil {
		log.Error().Err(err).Msg("handle error")
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("handle error")
		return
	}
	log.Info().Msg(string(body))
}
