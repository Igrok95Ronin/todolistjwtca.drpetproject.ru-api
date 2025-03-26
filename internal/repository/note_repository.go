package repository

import (
	"context"
	"database/sql"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"time"
)

// NoteRepository - интерфейс для работы с заметками
type NoteRepository interface {
	GetAllNotesFromDB(ctx context.Context, userID int64) ([]models.AllNotes, error)
	InsertNoteToDB(ctx context.Context, userID int64, note string, createdAt time.Time) error
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{
		db: db,
	}
}

// GetAllNotesFromDB - получаем все заметки из БД
func (r *noteRepository) GetAllNotesFromDB(ctx context.Context, userID int64) ([]models.AllNotes, error) {
	query := "SELECT id,note,completed,user_id,created_at FROM all_notes WHERE user_id = $1"

	// Используем QueryContext вместо QueryRowContext для множественных записей
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.AllNotes

	// Итерируемся по всем строкам
	for rows.Next() {
		var note models.AllNotes
		err = rows.Scan(
			&note.ID,
			&note.Note,
			&note.Completed,
			&note.UserID,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	// Проверяем ошибки после итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// InsertNoteToDB - добавить новую заметку в БД
func (r *noteRepository) InsertNoteToDB(ctx context.Context, userID int64, note string, createdAt time.Time) error {
	query := "INSERT INTO all_notes (note,user_id,created_at) VALUES ($1, $2, $3)"

	// Используйте ExecContext для операций INSERT/UPDATE/DELETE
	_, err := r.db.ExecContext(ctx, query, note, userID, createdAt)
	return err
}
