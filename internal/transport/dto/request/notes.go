package request

// CreateNote DTO для входящего запроса
type CreateNoteDTO struct {
	Note string `json:"note"`
}

// UpdateNote DTO для входящего запроса
type UpdateNoteDTO struct {
	CreateNoteDTO
}

// CheckNoteDTO DTO для входящего запроса
type CheckNoteDTO struct {
	Check bool `json:"check"`
}
