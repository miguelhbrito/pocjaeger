package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pocjaeger/tracing"
)

func main() {

	tracer, closer := tracing.InitJaeger("client-hello-world-tracing")
	defer closer.Close()

	span := tracer.StartSpan("say-hello-world-client")
	defer span.Finish()

	data := "Data to send to server"

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	_, err := DoRequest(ctx,data)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}
}