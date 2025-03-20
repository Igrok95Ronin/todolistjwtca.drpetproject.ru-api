package handlers

import (
	"encoding/json"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/service"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/httperror"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// UserHandler обрабатывает запросы, связанные с users
type UserHandler struct {
	service service.UserService
	logger  *logging.Logger
}

// NewUserHandler создаёт новый обработчик users
func NewUserHandler(service service.UserService, logger *logging.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// Регистрация (создание нового пользователя)
func (h *UserHandler) register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var users models.Users

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		httperror.WriteJSONError(w, "Ошибка декодирования в json", err, http.StatusBadRequest)
		h.logger.Errorf("Ошибка декодирования в json: %s", err)
		return
	}

	// UserExists проверяем есть ли пользователь и регистрирует нового пользователя
	if err := h.service.UserExists(ctx, users); err != nil {
		h.logger.Errorf("Ошибка при регистрации пользователя: %s", err)
		httperror.WriteJSONError(w, "Ошибка при регистрации пользователя", err, http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("Пользователь успешно зарегистрирован"))
	if err != nil {
		h.logger.Errorf("Обработка ошибки ответа: %s", err)
	}
}

// Логин (получение access и refresh токенов)
func (h *UserHandler) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var users models.Users

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		httperror.WriteJSONError(w, "Ошибка декодирования в json", err, http.StatusBadRequest)
		h.logger.Errorf("Ошибка декодирования в json: %s", err)
		return
	}

	if err := h.service.Login(ctx, w, users); err != nil {
		h.logger.Errorf("Ошибка при авторизации пользователя: %s", err)
		httperror.WriteJSONError(w, "Ошибка при авторизации пользователя", err, http.StatusInternalServerError)
		return
	}

	// Ответ для клиента
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Авторизация прошла успешно"))
	if err != nil {
		h.logger.Errorf("Ошибка авторизации: %s", err)
	}
}
