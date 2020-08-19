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

	span := tracer.StartSpan("say-hello-client")
	span.SetTag("client-span", "hi")
	defer span.Finish()

	data := "Data to seend to server"

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	_, err := Do(ctx,data)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}
}