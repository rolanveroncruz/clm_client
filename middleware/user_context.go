package middleware

import "context"

type middlewareKey int64

const (
	UserKey middlewareKey = iota
)

func ContextWithUser(ctx context.Context, userEmail string) context.Context {
	return context.WithValue(ctx, UserKey, userEmail)
}

func UserFromContext(ctx context.Context) interface{} {
	return ctx.Value(UserKey)
}
