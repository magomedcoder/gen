package domain

import "errors"

var (
	ErrUnauthorized = errors.New("недостаточно прав")
	ErrNoRunners    = errors.New("нет активных раннеров")
)
