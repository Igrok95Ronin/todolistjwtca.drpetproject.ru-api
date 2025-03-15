package repository

import (
	"database/sql"
	"fmt"

	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // Подключаем драйвер PostgreSQL
)

// NewDB создает подключение к БД
func NewDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Host, cfg.DB.Port, cfg.DB.SslMode, cfg.DB.TimeZone,
	)

	db, err := sql.Open("pgx", dsn) // Используем pgx вместо pq
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("не удалось выполнить ping к БД: %w", err)
	}

	return db, nil
}

// CloseDB закрывает соединение с БД
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
