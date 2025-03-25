package repository

import (
	"context"
	"database/sql"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
)

// UserRepository - интерфейс для работы с пользователями
type UserRepository interface {
	UserExists(userName, email string, ctx context.Context) error
	Register(users models.Users, ctx context.Context) error
	GetUser(ctx context.Context, users models.Users, email string) (*models.Users, error)
	UpdateRefreshToken(ctx context.Context, userID int64, refreshToken string) error
	FindUserByRefreshToken(ctx context.Context, users models.Users, userID int64) (*models.Users, error)
	GetUserProfileDB(ctx context.Context, userID int64) (*models.Users, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// UserExists проверяем есть ли пользователь в бд
func (r *userRepository) UserExists(userName, email string, ctx context.Context) error {
	query := "SELECT 1 FROM users WHERE user_name = $1 OR email = $2 LIMIT 1"
	var exists int
	return r.db.QueryRowContext(ctx, query, userName, email).Scan(&exists)
}

// Register Сохраняем пользователя в бд
func (r *userRepository) Register(users models.Users, ctx context.Context) error {
	query := "INSERT INTO users (user_name, email, password_hash, created_at) VALUES ($1,$2, $3, NOW())"
	_, err := r.db.ExecContext(ctx, query, users.UserName, users.Email, users.PasswordHash)
	return err
}

// GetUser получаем пользователя из БД
func (r *userRepository) GetUser(ctx context.Context, users models.Users, email string) (*models.Users, error) {
	query := "SELECT id, email, refresh_token, password_hash FROM users WHERE email = $1 LIMIT 1"

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&users.ID,
		&users.Email,
		&users.RefreshToken,
		&users.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

// UpdateRefreshToken обновляем refreshToken в БД
func (r *userRepository) UpdateRefreshToken(ctx context.Context, userID int64, refreshToken string) error {
	query := "UPDATE users SET refresh_token = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, refreshToken, userID)
	return err
}

// FindUserByRefreshToken Проверим, что refresh-токен совпадает с тем, что хранится в базе
func (r *userRepository) FindUserByRefreshToken(ctx context.Context, users models.Users, userID int64) (*models.Users, error) {
	query := "SELECT id, refresh_token FROM users WHERE id = $1 LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&users.ID,
		&users.RefreshToken,
	)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

// GetUserProfile Получить данные о текущем пользователе из БД
func (r *userRepository) GetUserProfileDB(ctx context.Context, userID int64) (*models.Users, error) {
	query := "SELECT id, user_name, email FROM users WHERE id = $1 LIMIT 1"

	var user models.Users

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
