package middleware

import (
	"context"
	"net/http"
)

type contextKey struct {
	uuid string
}

var HttpWriterKey = &contextKey{"httpWriter"}
var HttpReaderKey = &contextKey{"httpReader"}

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// put it in context
		ctx := context.WithValue(r.Context(), HttpWriterKey, w)
		ctx = context.WithValue(ctx, HttpReaderKey, w)
		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
