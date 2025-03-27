package service

import (
	"context"
	"fmt"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/errors"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/repository"
	"html"
	"strings"
	"time"
	"unicode/utf8"
)

// NoteService - интерфейс для работы с бизнес-логикой заметок
type NoteService interface {
	GetAllNotes(ctx context.Context, userID int64) ([]models.AllNotes, error)
	ValidateNoteBeforeInserting(ctx context.Context, userID int64, note string) error
	UpdateNoteDataValidation(ctx context.Context, id int64, note string) error
}

type noteService struct {
	repo repository.NoteRepository
	cfg  *config.Config
}

func NewNoteService(repo repository.NoteRepository, cfg *config.Config) NoteService {
	return &noteService{
		repo: repo,
		cfg:  cfg,
	}
}

// GetAllNotes - получаем все заметки
func (s *noteService) GetAllNotes(ctx context.Context, userID int64) ([]models.AllNotes, error) {
	allNotesFromDB, err := s.repo.GetAllNotesFromDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	return allNotesFromDB, err
}

// ValidateTheNoteBeforeInserting - валидация заметки перед вставкой
func (s *noteService) ValidateNoteBeforeInserting(ctx context.Context, userID int64, note string) error {

	note = html.EscapeString(strings.TrimSpace(note))

	if utf8.RuneCountInString(note) < 3 {
		return errors.ErrNoteTooShort
	}

	createdAt := time.Now().UTC() // UTC для универсальности

	// InsertNoteToDB - добавить новую заметку в БД
	if err := s.repo.InsertNoteToDB(ctx, userID, note, createdAt); err != nil {
		return fmt.Errorf("%w: %w", errors.ErrNoteFailed, err)
	}

	return nil
}

// UpdateNoteDataValidation - обновление заметки, валидация данных
func (s *noteService) UpdateNoteDataValidation(ctx context.Context, id int64, note string) error {
	note = html.EscapeString(strings.TrimSpace(note))

	if utf8.RuneCountInString(note) < 3 {
		return errors.ErrNoteTooShort
	}

	if id <= 0 {
		return errors.ErrIDCannotBeNegativeOrEqualToZero
	}

	if err := s.repo.UpdateNoteToDB(ctx, id, note); err != nil {
		return err
	}

	return nil
}
