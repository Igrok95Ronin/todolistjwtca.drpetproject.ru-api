package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// UserService - интерфейс для работы с бизнес-логикой заметок
type UserService interface {
	UserExists(ctx context.Context, users models.Users) error
	Login(ctx context.Context, w http.ResponseWriter, users models.Users) error
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

// UserExists проверяем есть ли пользователь регистрируем нового пользователя
func (s *userService) UserExists(ctx context.Context, users models.Users) error {

	userName := strings.TrimSpace(users.UserName)
	email := strings.TrimSpace(users.Email)
	password := strings.TrimSpace(users.PasswordHash)

	// Проверка, что поля заполнены
	if userName == "" || email == "" || password == "" {
		return fmt.Errorf("Все поля (username, email, password) обязательны")
	}

	//Запрет на выполнение скриптов
	userName = template.HTMLEscapeString(userName)
	email = template.HTMLEscapeString(email)
	password = template.HTMLEscapeString(password)

	// Проверка валидности email
	if err := ValidateEmail(email); err != nil {
		return fmt.Errorf("Неверный формат email: %s", email)
	}

	// UserExists проверяем есть ли пользователь в бд
	err := s.repo.UserExists(userName, email, ctx)
	if err == nil { // Если ошибки нет, значит пользователь найден
		return errors.New("Пользователь с таким username или email уже существует")
	}

	if err != sql.ErrNoRows { // Если ошибка не sql.ErrNoRows, значит это другая проблема
		return fmt.Errorf("ошибка при проверке пользователя: %w", err)
	}

	// Хешируем пароль
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("Ошибка при хешировании пароля: %w", err)
	}

	// Создаём объект нового пользователя
	newUser := models.Users{
		UserName:     userName,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// Сохраняем пользователя
	if err = s.repo.Register(newUser, ctx); err != nil {
		return fmt.Errorf("Ошибка при сохранении пользователя: %s", err)
	}

	return nil
}

// Login проверяем есть ли пользователь (получение access и refresh токенов)
func (s *userService) Login(ctx context.Context, w http.ResponseWriter, users models.Users) error {

	email := strings.TrimSpace(users.Email)
	password := strings.TrimSpace(users.PasswordHash)

	// Проверка, что поля заполнены
	if email == "" || password == "" {
		return fmt.Errorf("Все поля (email, password) обязательны")
	}

	//Запрет на выполнение скриптов
	email = template.HTMLEscapeString(email)
	password = template.HTMLEscapeString(password)

	// Проверка валидности email
	if err := ValidateEmail(email); err != nil {
		return fmt.Errorf("Неверный формат email: %s", email)
	}

	// UserExists проверяем есть ли пользователь в бд
	user, err := s.repo.GetUser(ctx, users, email)

	if err != nil {
		return fmt.Errorf("ошибка при проверке пользователя: %w", err)
	}
	if user == nil {
		return errors.New("Неверный email или пароль")
	}

	// Проверяем пароль (сравниваем с хешем в базе)
	if !CheckPasswordHash(password, user.PasswordHash) {
		return fmt.Errorf("Неверный email или пароль: %w", err)
	}

	// Генерируем access-токен
	accessToken, err := GenerateAccessToken(s, user.ID)
	if err != nil {
		return fmt.Errorf("Ошибка при генерации access-токена: %w", err)
	}

	// Генерируем refresh-токен
	refreshToken, err := GenerateRefreshToken(s, user.ID)
	if err != nil {
		return fmt.Errorf("Ошибка при генерации refresh-токена: %w", err)
	}

	// Сохраняем refresh-токен у пользователя в базе (на практике лучше хранить хеш)
	//user.RefreshToken = refreshToken
	if err = s.repo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return fmt.Errorf("ошибка при сохранении refresh-токена: %w", err)
	}

	// Устанавливаем access-токен в куки (жизнь 15 минут)
	// HttpOnly: true означает, что кука не доступна из JavaScript (защита от XSS).
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Path:     "/",
		// Secure:   true, // Использовать при HTTPS
		// SameSite: http.SameSiteStrictMode,
	})

	// Устанавливаем refresh-токен в куки (жизнь 30 дней)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		// Secure:   true, // Использовать при HTTPS
		// SameSite: http.SameSiteStrictMode,
	})

	return nil
}

//---------------------------------------------------------------------------------------
//                                 УТИЛИТНЫЕ ФУНКЦИИ
//---------------------------------------------------------------------------------------

// HashPassword - хеширует пароль с помощью bcrypt (с cost = bcrypt.DefaultCost).
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword вернёт хеш пароля.
	// bcrypt.DefaultCost по умолчанию равен 10 (можно увеличить, чтобы усложнить подбор).
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка валидности email
func ValidateEmail(email string) error {
	// Проверка длины email
	if len(email) < 4 {
		return fmt.Errorf("email должен быть не менее 4 символов")
	}

	// Проверка наличия символа "@" в email
	if !strings.Contains(email, "@") {
		return fmt.Errorf("email должен содержать символ '@'")
	}

	// Проверка позиции символа "@" (не должен быть первым или последним символом)
	if strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return fmt.Errorf("email не может начинаться или заканчиваться на '@'")
	}

	// Дополнительно: базовая проверка формата email с помощью регулярного выражения
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return fmt.Errorf("ошибка проверки email: %v", err)
	}
	if !matched {
		return fmt.Errorf("email не соответствует формату")
	}

	return nil
}

// CheckPasswordHash - проверяет соответствие "сырого" пароля и хеша.
func CheckPasswordHash(password, hash string) bool {
	// bcrypt.CompareHashAndPassword вернёт nil, если всё совпадает.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAccessToken - генерирует access-токен с временем жизни 15 минут.
// Внутри указываем UserID и стандартные поля (ExpiresAt, IssuedAt, NotBefore).
func GenerateAccessToken(s *userService, userID int64) (string, error) {
	// Создаём claims.
	claims := models.MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Токен протухнет через 15 минут
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       // Время выпуска
			NotBefore: jwt.NewNumericDate(time.Now()),                       // Не раньше этого времени
		},
	}

	// Создаём токен с алгоритмом HS256.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен нашим секретным ключом.
	return token.SignedString([]byte(s.cfg.Token.Access))
}

// GenerateRefreshToken - генерирует refresh-токен с временем жизни 30 дней.
func GenerateRefreshToken(s *userService, userID int64) (string, error) {
	claims := models.MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 дней
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Token.Refresh))
}
