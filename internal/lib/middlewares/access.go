package middlewares

import (
	"net/http"

	"github.com/prabhatlabs/go-mail-server/internal/lib/env"
	"github.com/prabhatlabs/go-mail-server/internal/lib/response"
)

func AccessCodeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac := r.Header.Get("X-Access-Code")
		if ac != env.Vars.ACCESS_CODE {
			response.Unauthorized(w, "Authentication required")
			return
		}
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
