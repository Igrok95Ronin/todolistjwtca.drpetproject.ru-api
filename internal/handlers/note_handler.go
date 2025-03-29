package handlers

import (
	"encoding/json"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/errors"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/service"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/internal/transport/dto/request"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/httperror"
	"github.com/Igrok95Ronin/todolistjwtca.drpetproject.ru-api.git/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
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
		httperror.WriteJSONError(w, errors.ErrJSONNewDecoder.Error(), err, http.StatusBadRequest)
		h.logger.Errorf("%s: %s", errors.ErrJSONNewDecoder, err)
		return
	}

	// ValidateTheNoteBeforeInserting - валидация заметки перед вставкой
	if err := h.noteService.ValidateNoteBeforeInserting(ctx, userID, req.Note); err != nil {
		httperror.WriteJSONError(w, "Ошибка при добавлении новой заметки", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Обновить заметку
func (h *NoteHandler) updateNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	var req request.UpdateNoteDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Если произошла ошибка декодирования, возвращаем клиенту ошибку с кодом 400
		httperror.WriteJSONError(w, errors.ErrJSONNewDecoder.Error(), err, http.StatusBadRequest)
		h.logger.Errorf("%s: %s", errors.ErrJSONNewDecoder, err)
		return
	}

	id, _ := strconv.Atoi(ps.ByName("id"))

	// UpdateNoteDataValidation - обновление заметки, валидация данных
	if err := h.noteService.UpdateNoteDataValidation(ctx, int64(id), req.Note); err != nil {
		httperror.WriteJSONError(w, "Ошибка при обновления записи в БД", err, http.StatusInternalServerError)
		h.logger.Errorf("Ошибка при обновлении записи по id: %v %s", id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Удалить конкретную заметку
func (h *NoteHandler) deleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id, _ := strconv.Atoi(ps.ByName("id"))

	if err := h.noteService.DeleteNote(ctx, int64(id)); err != nil {
		httperror.WriteJSONError(w, errors.ErrDeleteNote.Error(), err, http.StatusInternalServerError)
		h.logger.Errorf("%s : %v : %s", errors.ErrDeleteNote, id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Отметить заметку выполненной
func (h *NoteHandler) markNoteCompleted(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id, _ := strconv.Atoi(ps.ByName("id"))

	var req request.CheckNoteDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Если произошла ошибка декодирования, возвращаем клиенту ошибку с кодом 400
		httperror.WriteJSONError(w, errors.ErrJSONNewDecoder.Error(), err, http.StatusBadRequest)
		h.logger.Errorf("%s: %s", errors.ErrJSONNewDecoder, err)
		return
	}

	// MarkNoteCompleted - Отметить заметку выполненной, валидация данных
	if err := h.noteService.MarkNoteCompleted(ctx, int64(id), req.Check); err != nil {
		httperror.WriteJSONError(w, errors.ErrNoteToUpdate.Error(), err, http.StatusInternalServerError)
		h.logger.Errorf("%s : %v : %s", errors.ErrNoteToUpdate, id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
