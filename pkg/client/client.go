package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

type Response struct {
	TraceID string `json:"trace-id"`
	SpanID  string `json:"span-id,omitempty"`
}

func DoRequest(ctx context.Context) (*http.Response, error) {
	resp := &http.Response{}
	//span, _ := opentracing.StartSpanFromContext(ctx, "Doing request")
	//defer span.Finish()
	//tid := tracing.GetTraceID(span)

	url := "http://localhost:8080/server-two"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp, err
	}

	span := opentracing.SpanFromContext(ctx)

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

