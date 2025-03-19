package main

import (
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/handlers"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/middleware"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/repository"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func main() {

	// Загружаем конфигурацию
	cfg := config.GetConfig()

	// Настраиваем логгер
	logger := logging.GetLogger()

	// Инициализируем базу данных (в слое repository)
	db, err := repository.NewDB(cfg)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer repository.CloseDB(db)

	// Создаем роутер
	router := httprouter.New()

	// Инициализируем обработчики (handlers) и передаем им зависимости
	handler := handlers.NewHandler(cfg, logger, db)
	handler.RegisterRoutes(router)

	// Обработка cors, Context
	corsHandler := middleware.CorsSettings().Handler(middleware.RequestContext(router))

	// Запускаем сервер
	start(corsHandler, cfg, logger)
}

func start(router http.Handler, cfg *config.Config, logger *logging.Logger) {
	const timeout = 15 * time.Second

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		IdleTimeout:  timeout,
	}

	logger.Infof("Сервер запущен на %v", cfg.Port)
	logger.Fatal(server.ListenAndServe())
}
