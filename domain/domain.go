package domain

import "time"

type AIChatMessageRole string

const (
	AIChatMessageRoleSystem    AIChatMessageRole = "system"
	AIChatMessageRoleUser      AIChatMessageRole = "user"
	AIChatMessageRoleAssistant AIChatMessageRole = "assistant"
)

type AIChatMessage struct {
	Id        int64
	SessionId int64
	Content   string
	Role      AIChatMessageRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AIFromProtoRole(role string) AIChatMessageRole {
	switch role {
	case "system":
		return AIChatMessageRoleSystem
	case "user":
		return AIChatMessageRoleUser
	case "assistant":
		return AIChatMessageRoleAssistant
	default:
		return AIChatMessageRoleUser
	}
}
