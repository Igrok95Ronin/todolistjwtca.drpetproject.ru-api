package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/service"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/httperror"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
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

// RefreshHandler - обработчик обновления токенов.
func (h *UserHandler) refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var users models.Users
	// 1. Извлекаем refresh_token из куки
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		httperror.WriteJSONError(w, "Необходим refresh_token (cookie отсутствует)", err, http.StatusUnauthorized)
		h.logger.Errorf("Необходим refresh_token (cookie отсутствует): %s", err)
		return
	}
	refreshToken := refreshCookie.Value

	// 2. Валидируем refresh-токен
	if err = h.service.Refresh(ctx, w, users, refreshToken); err != nil {
		httperror.WriteJSONError(w, "Невалидный или просроченный refresh-токен", err, http.StatusUnauthorized)
		h.logger.Errorf("Невалидный или просроченный refresh-токен: %s", err)
		return
	}
}

// ProtectedHandler - обработчик примера защищённого маршрута.
// Доступ сюда возможен только через Auth.
func (h *UserHandler) protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//// Достаём user_id из контекста (установили в AuthMiddleware).
	//userID, ok := r.Context().Value("user_id").(uint)
	//if !ok {
	//	// Если что-то пошло не так и user_id не смогли получить
	//	http.Error(w, "Не удалось получить user_id из контекста", http.StatusInternalServerError)
	//	return
	//}

	// Если всё ок, возвращаем сообщение, что доступ разрешён.
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("Доступ к защищённому маршруту разрешен.")))
	if err != nil {
		h.logger.Error(err)
	}
}

// Выход из системы
func (h *UserHandler) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Устанавливаем куки с прошедшей датой
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // просрочен
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Вы успешно вышли из системы"))
	if err != nil {
		h.logger.Error(err)
	}
}

type UserProfileResponse struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

// Получить данные о текущем пользователе
func (h *UserHandler) getUserProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.logger.Error("Не удалось получить user_id из контекста")
		httperror.WriteJSONError(w, "Не удалось получить user_id из контекста", nil, http.StatusUnauthorized)
		return
	}

	userProfile, err := h.service.GetUserProfile(ctx, userID)
	if err != nil {
		httperror.WriteJSONError(w, "Возможно данные о пользователе отсутствуют", err, http.StatusBadRequest)
		h.logger.Errorf("Возможно данные о пользователе отсутствуют: %s", err)
		return
	}

	userProfileResponse := UserProfileResponse{
		ID:       userProfile.ID,
		UserName: userProfile.UserName,
		Email:    userProfile.Email,
	}

	// Отправляем JSON-ответ с user_name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userProfileResponse)
	if err != nil {
		h.logger.Error(err)
	}
}
