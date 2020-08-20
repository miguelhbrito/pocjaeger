package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"io/ioutil"
	"net/http"
	"net/url"
)

func DoRequest(ctx context.Context, data string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "doing-request")
	span.SetTag("spanID", span.Context().(jaeger.SpanContext).SpanID())
	span.SetTag("traceID", span.Context().(jaeger.SpanContext).TraceID())
	defer span.Finish()

	v := url.Values{}
	v.Set("hello-world", data)
	url := "http://localhost:8000/helloWorld?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}
