package request

// CreateNote DTO для входящего запроса
type CreateNoteDTO struct {
	Note string `json:"note"`
}

// UpdateNote DTO для входящего запроса
type UpdateNoteDTO struct {
	CreateNoteDTO
}
