package service

import "github.com/magomedcoder/llm-runner/domain"

func ChatRoleString(role domain.AIChatMessageRole) string {
	switch role {
	case domain.AIChatMessageRoleSystem:
		return "system"
	case domain.AIChatMessageRoleAssistant:
		return "assistant"
	default:
		return "user"
	}
}
