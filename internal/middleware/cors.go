package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

// Обработка CORS
func CorsSettings() *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodPut,
			http.MethodOptions, // Добавлен OPTIONS для preflight-запросов
		},
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:81", "https://todolistjwt.drpetproject.ru"},
		AllowCredentials: true, // Разрешаем отправку cookie (credentials)
		AllowedHeaders: []string{
			"X-Api-Password",
			"Content-Type",
			"Authorization",
			"X-Requested-With", // Добавлен заголовок из corsMiddleware
		},
		OptionsPassthrough: false, // Прекращаем обработку preflight-запросов после CORS
	})
}
