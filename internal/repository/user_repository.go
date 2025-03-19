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
