package domain

import "errors"

var (
	ErrUnauthorized                = errors.New("недостаточно прав")
	ErrNoRunners                   = errors.New("нет активных раннеров")
	ErrRegenerateToolsNotSupported = errors.New("перегенерация недоступна при включённых инструментах")
)
