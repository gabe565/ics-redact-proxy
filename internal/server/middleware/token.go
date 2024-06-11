package middleware

import (
	"net/http"
	"slices"
)

func Token(tokens ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if len(tokens) == 0 {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(tokens, r.URL.Query().Get("token")) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
		})
	}
}
