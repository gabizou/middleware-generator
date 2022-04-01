package example

import (
	"context"

	"example/domain"
)

//go:generate cd .. && go run ./cmd/generator Service SvcMiddleware tracer

type Service interface {
	Foo(ctx context.Context, bar string) domain.Foo
}

type SvcMiddleware func(Service) Service

type Pill struct {
}

type unexported []map[string]*[]*interface{}

//go:generate cd .. && go run ./cmd/generator Repository RepoMiddleware tracer

type Repository interface {
	Find(ctx context.Context, id string) (*domain.Foo, error)
	Foo(ctx context.Context) (anInt int, aBool bool, aSlice []*domain.Foo, complexSlice []*[]interface{}, aMap map[string]*interface{})
	Bar(ctx context.Context, astruct struct{ name string }) **interface {
		aFunc(inner func(ctx context.Context, uint2 uint) (string, error, unexported))
	}

	Baz(ctx context.Context) func(ctx context.Context) error
}

type RepoMiddleware func(Repository) Repository
