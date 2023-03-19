package stub

import "context"

type daoStubCtxKeyType int

const (
	helloCtxKey daoStubCtxKeyType = iota
)

// WithHelloDao adds hello DAO stub to the context
func WithHelloDao(parent context.Context) context.Context {
	if parent.Value(helloCtxKey) != nil {
		panic("dao already in the context")
	}

	ctx := context.WithValue(parent, helloCtxKey, &helloDaoStub{})
	return ctx
}

func getHelloDaoStub(ctx context.Context) *helloDaoStub {
	var ok bool
	var resDao *helloDaoStub
	if resDao, ok = ctx.Value(helloCtxKey).(*helloDaoStub); !ok {
		panic("dao not in context")
	}
	return resDao
}
