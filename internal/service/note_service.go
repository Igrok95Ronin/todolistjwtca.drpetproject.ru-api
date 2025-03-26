package service

import (
	"context"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/config"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/models"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/repository"
)

// NoteService - интерфейс для работы с бизнес-логикой заметок
type NoteService interface {
	GetAllNotes(ctx context.Context, userID int64) ([]models.AllNotes, error)
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
