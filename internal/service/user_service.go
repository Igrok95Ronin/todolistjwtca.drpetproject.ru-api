package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"regexp"
	"strings"
)

// UserService - интерфейс для работы с бизнес-логикой заметок
type UserService interface {
	UserExists(users models.Users, ctx context.Context) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// UserExists проверяем есть ли пользователь
func (s *userService) UserExists(users models.Users, ctx context.Context) error {

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
