package middleware

import (
	"context"
	"net/http"

	"stars-server/app/models"
)

func ContextEnricher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add real game id reading
		r = r.WithContext(context.WithValue(r.Context(), models.GameInfoKey{}, 1))

		next.ServeHTTP(w, r)
	})
}
