package openai

import "context"

type sessionContextKey struct{}

func ContextSetSession(ctx context.Context, sesh *Chat) context.Context {
	return context.WithValue(ctx, sessionContextKey{}, sesh)
}

func ContextGetSession(ctx context.Context) *Chat {
	user, ok := ctx.Value(sessionContextKey{}).(*Chat)
	if !ok {
		panic("missing session in context")
	}

	return user
}
