package handlers

import (
	"database/sql"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/middleware"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/repository"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/service"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

// Handler управляет роутами
type Handler struct {
	cfg      *config.Config
	logger   *logging.Logger
	userRepo repository.UserRepository
	userSvc  service.UserService
	noteRepo repository.NoteRepository
	noteSvc  service.NoteService
}

// NewHandler создаёт новый обработчик
func NewHandler(cfg *config.Config, logger *logging.Logger, db *sql.DB) *Handler {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, cfg)

	noteRepo := repository.NewNoteRepository(db)
	noteSvc := service.NewNoteService(noteRepo, cfg)

	return &Handler{
		cfg:      cfg,
		logger:   logger,
		userRepo: userRepo,
		userSvc:  userSvc,
		noteRepo: noteRepo,
		noteSvc:  noteSvc,
	}
}

// RegisterRoutes регистрирует маршруты
func (h *Handler) RegisterRoutes(router *httprouter.Router) {
	userHandler := NewUserHandler(h.userSvc, h.logger)
	noteHandler := NewNoteHandler(h.noteSvc, h.logger)

	router.POST("/register", userHandler.register)                       // Регистрация (создание нового пользователя)
	router.POST("/login", userHandler.login)                             // Логин (получение access и refresh токенов)
	router.POST("/refresh", userHandler.refresh)                         // Обновление (refresh) токенов
	router.POST("/logout", userHandler.logout)                           // Выход из системы
	router.GET("/protected", middleware.Auth(userHandler.protected))     // Защищённый маршрут, доступный только при наличии валидного access-токена
	router.GET("/users/me", middleware.Auth(userHandler.getUserProfile)) // Получить данные о текущем пользователе

	router.GET("/notes", middleware.Auth(noteHandler.getAllNotes))                     // Получить все заметки
	router.POST("/notes", middleware.Auth(noteHandler.createPost))                     // Создать заметку
	router.PUT("/notes/:id", middleware.Auth(noteHandler.updateNote))                  // Обновить заметку
	router.DELETE("/notes/:id", middleware.Auth(noteHandler.deleteNote))               // Удалить конкретную заметку
	router.PUT("/notes/:id/completed", middleware.Auth(noteHandler.markNoteCompleted)) // Отметить заметку выполненной
	router.DELETE("/notes", middleware.Auth(noteHandler.deleteAllNotes))               // Удалить все заметки
}
