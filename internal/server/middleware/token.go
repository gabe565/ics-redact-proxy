package middleware

import (
	"net/http"
	"slices"
)

func Token(h http.Handler, tokens ...string) http.Handler {
	if len(tokens) == 0 {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(tokens, r.URL.Query().Get("token")) {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}
