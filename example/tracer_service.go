// Code generated by "logger-middleware Service SvcMiddleware tracer"; DO NOT EDIT.

package example

import (
	"context"
	domain "example/domain"

	zipkingo "github.com/openzipkin/zipkin-go"
)

func NewServiceTracer(tracer zipkingo.Tracer) SvcMiddleware {
	return func(s Service) Service {
		return &tracerS{
			s:  s,
			tr: tracer,
		}
	}
}

type tracerS struct {
	tr zipkingo.Tracer
	s  Service
}

func (t *tracerS) Foo(ctx context.Context, bar string) domain.Foo {
	span, ctx := t.tr.StartSpanFromContext(ctx, "Foo")

	defer func() {
		span.Finish()
	}()

	return t.s.Foo(ctx, bar)
}