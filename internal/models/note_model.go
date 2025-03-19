package models

import "time"

// Структура для таблицы all_notes
type AllNotes struct {
	ID        int64     `json:"ID" gorm:"primaryKey;column:id"`    // Первичный ключ
	Note      string    `json:"note" gorm:"column:note"`           // Поле заметки
	Completed bool      `json:"completed" gorm:"column:completed"` // Статус выполнения
	UserID    int64     `json:"userID" gorm:"column:user_id"`      // Связь с таблицей users
	CreatedAt time.Time `gorm:"column:created_at"`                 // Дата создания
}
