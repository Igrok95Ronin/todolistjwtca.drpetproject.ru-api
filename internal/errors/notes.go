package errors

import "errors"

var (
	ErrNoteTooShort = errors.New("Слишком короткая заметка")
	ErrNoteFailed   = errors.New("Вставить заметку не удалось")
	ErrNoteNotFound = errors.New("Заметка Не Найдена")

	ErrIDCannotBeNegativeOrEqualToZero = errors.New("ID не может быть отрицательным или равным 0")

	ErrNoteToUpdate           = errors.New("Не удалось обновить заметку")
	FailedToCheckAffectedRows = errors.New("Не удалось проверить затронутые строки")
)
