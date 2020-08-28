package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"reflect"
	"runtime"
)

func RunTracedFunction(fn interface{}, params ...interface{}) (result []reflect.Value) {
	vf := reflect.ValueOf(fn)
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		inputs[k] = reflect.ValueOf(in)
	}

	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	ctx := params[0].(context.Context)
	parentSpan := opentracing.SpanFromContext(ctx)
	sp := opentracing.StartSpan("Function - " + name, opentracing.ChildOf(parentSpan.Context()))
	sp.SetTag("function", name)

	ctx = opentracing.ContextWithSpan(ctx, sp)
	inputs[0] = reflect.ValueOf(ctx)

	result = vf.Call(inputs)

	sp.Finish()
	return
}

func MakeTracedFunction(fn interface{}) interface{} {
	vf := reflect.ValueOf(fn)
	wrapperF := reflect.MakeFunc(reflect.TypeOf(fn),
		func(in []reflect.Value) []reflect.Value {
			name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
			ctx := in[0].Interface().(context.Context)
			parentSpan := opentracing.SpanFromContext(ctx)
			sp := opentracing.StartSpan("Function - " + name, opentracing.ChildOf(parentSpan.Context()))
			sp.SetTag("function", name)

			out := vf.Call(in)

			sp.Finish()
			return out
		},
	)
	return wrapperF.Interface()
}