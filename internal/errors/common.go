package errors

import "errors"

var (
	ErrJSONNewDecoder               = errors.New("Ошибка декодирования в JSON")
	ErrFailedToGetUserIDFromContext = errors.New("Не удалось получить user_id из контекста")
)
