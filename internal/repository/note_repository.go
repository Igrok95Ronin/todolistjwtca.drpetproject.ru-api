package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/errors"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"time"
)

// NoteRepository - интерфейс для работы с заметками
type NoteRepository interface {
	GetAllNotesFromDB(ctx context.Context, userID int64) ([]models.AllNotes, error)
	InsertNoteToDB(ctx context.Context, userID int64, note string, createdAt time.Time) error
	UpdateNoteToDB(ctx context.Context, id int64, note string) error
	DeleteNoteFromDB(ctx context.Context, id int64) error
	MarkNoteCompletedToDB(ctx context.Context, id int64, check bool) error
	DeleteAllNotesFromDB(ctx context.Context, userID int64) error
	DeleteAllCompletedNotesFromDB(ctx context.Context, userID int64) error
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

// UpdateNoteToDB - обновить заметку в БД
func (r *noteRepository) UpdateNoteToDB(ctx context.Context, id int64, note string) error {
	query := "UPDATE all_notes SET note = $1 WHERE id = $2"

	// Используйте ExecContext для операций INSERT/UPDATE/DELETE
	result, err := r.db.ExecContext(ctx, query, note, id)
	if err != nil {
		return errors.ErrNoteToUpdate
	}

	// Проверяем, что запись была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.FailedToCheckAffectedRows
	}

	if rowsAffected == 0 {
		return errors.ErrNoteNotFound
	}

	return nil
}

// DeleteNoteFromDB - удалить заметку из БД
func (r *noteRepository) DeleteNoteFromDB(ctx context.Context, id int64) error {
	query := "DELETE FROM all_notes WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDeleteNote, err)
	}

	// Возвращает количество строк, затронутых обновлением, вставкой или удалением.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.FailedToCheckAffectedRows, err)
	}

	if rowsAffected == 0 {
		return errors.ErrNoteNotFound
	}

	return nil
}

// MarkNoteCompleted - Отметить заметку выполненной в БД
func (r *noteRepository) MarkNoteCompletedToDB(ctx context.Context, id int64, check bool) error {
	query := "UPDATE all_notes SET completed = $1 WHERE id = $2"

	// Используйте ExecContext для операций INSERT/UPDATE/DELETE
	result, err := r.db.ExecContext(ctx, query, check, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNoteToUpdate, err)
	}

	// Проверяем, что запись была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.FailedToCheckAffectedRows, err)
	}

	if rowsAffected == 0 {
		return errors.ErrNoteNotFound
	}

	return nil
}

// DeleteAllNotes - Удалить все заметки из БД
func (r *noteRepository) DeleteAllNotesFromDB(ctx context.Context, userID int64) error {
	query := "DELETE FROM all_notes WHERE user_id = $1"

	// Используйте ExecContext для операций INSERT/UPDATE/DELETE
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDeletingAllNotes, err)
	}

	// Проверяем, что запись была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.FailedToCheckAffectedRows, err)
	}

	if rowsAffected == 0 {
		return errors.ErrNoteNotFound
	}

	return nil
}

// DeleteAllCompletedNotesFromDB - Удалить все выполненные заметки из БД
func (r *noteRepository) DeleteAllCompletedNotesFromDB(ctx context.Context, userID int64) error {
	query := "DELETE FROM all_notes WHERE user_id = $1 AND completed = $2"

	// Используйте ExecContext для операций INSERT/UPDATE/DELETE
	result, err := r.db.ExecContext(ctx, query, userID, true)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrDeletingAllNotes, err)
	}

	// Проверяем, что запись была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.FailedToCheckAffectedRows, err)
	}

	if rowsAffected == 0 {
		return errors.ErrNoteNotFound
	}

	return nil
}
