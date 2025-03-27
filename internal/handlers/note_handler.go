package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/service"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/transport/rest/dto/request"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/httperror"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	ErrNoteJSONNewDecoder = errors.New("Ошибка декодирования в JSON")
)

// NoteHandler обрабатывает запросы, связанные с заметками
type NoteHandler struct {
	noteService service.NoteService
	logger      *logging.Logger
}

// NewNoteHandler создаёт новый обработчик заметок
func NewNoteHandler(noteService service.NoteService, logger *logging.Logger) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
		logger:      logger,
	}
}

// Получить все заметки
func (h *NoteHandler) getAllNotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.logger.Error("Не удалось получить user_id из контекста")
		httperror.WriteJSONError(w, "Не удалось получить user_id из контекста", nil, http.StatusInternalServerError)
		return
	}

	// GetAllNotes - получаем все заметки
	allNotes, err := h.noteService.GetAllNotes(ctx, userID)
	if err != nil {
		h.logger.Errorf("Ошибка при получения всех заметок", err)
		httperror.WriteJSONError(w, "Ошибка при получения всех заметок", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(allNotes); err != nil {
		h.logger.Errorf("Ошибка при отправке заметок на клиент", err)
		httperror.WriteJSONError(w, "Ошибка при отправке заметок на клиент", err, http.StatusInternalServerError)
	}
}

// Создать заметку
func (h *NoteHandler) createPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.logger.Error("Не удалось получить user_id из контекста")
		httperror.WriteJSONError(w, "Не удалось получить user_id", nil, http.StatusInternalServerError)
		return
	}

	var req request.CreateNoteDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Если произошла ошибка декодирования, возвращаем клиенту ошибку с кодом 400
		httperror.WriteJSONError(w, ErrNoteJSONNewDecoder.Error(), err, http.StatusBadRequest)
		// Логируем ошибку
		h.logger.Errorf("%s: %s", ErrNoteJSONNewDecoder, err)
		return
	}

	// ValidateTheNoteBeforeInserting - валидация заметки перед вставкой
	if err := h.noteService.ValidateNoteBeforeInserting(ctx, userID, req.Note); err != nil {
		httperror.WriteJSONError(w, "Ошибка при добавлении новой заметки", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
