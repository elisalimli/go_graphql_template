package middleware

import (
	"context"
	"net/http"

	myContext "github.com/elisalimli/go_graphql_template/context"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// put it in context
		ctx := context.WithValue(r.Context(), myContext.HttpReaderKey, r)
		ctx = context.WithValue(ctx, myContext.HttpWriterKey, w)

		refreshToken, err := r.Cookie("refresh_token")
		if err == nil {
			ctx = context.WithValue(ctx, myContext.CookieRefreshTokenKey, refreshToken)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
