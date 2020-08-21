package client

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"io/ioutil"
	"net/http"
)

func DoRequest(ctx context.Context) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "doing-request")
	defer span.Finish()

	sid := span.Context().(jaeger.SpanContext).SpanID()
	tid := span.Context().(jaeger.SpanContext).TraceID()
	span.SetTag("spanID", sid)
	span.SetTag("traceID", tid)

	tids := fmt.Sprintf("%s", tid)

	url := "http://localhost:8080/server-two"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("TraceID", tids)
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
