package middleware

import (
	"context"
	"net/http"
	"time"
)

// Middleware для установки контекста с тайм-аутом
func RequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем контекст с тайм-аутом 10 секунд
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel() // Отмена контекста после завершения

		// Передаем запрос с контекстом дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
