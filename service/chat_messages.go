package service

import (
	"github.com/magomedcoder/llm-runner/domain"
	"strings"
)

func NormalizeChatMessages(messages []*domain.AIChatMessage) []*domain.AIChatMessage {
	if len(messages) == 0 {
		return nil
	}

	var systemParts []string
	var rest []*domain.AIChatMessage

	for _, m := range messages {
		if m == nil {
			continue
		}

		if strings.TrimSpace(m.Content) == "" {
			continue
		}

		switch m.Role {
		case domain.AIChatMessageRoleSystem:
			systemParts = append(systemParts, strings.TrimSpace(m.Content))
		default:
			rest = append(rest, m)
		}
	}

	var out []*domain.AIChatMessage
	if len(systemParts) > 0 {
		merged := strings.Join(systemParts, "\n\n")
		out = append(out, domain.NewAIChatMessage(0, merged, domain.AIChatMessageRoleSystem))
	}
	out = append(out, rest...)

	return out
}

func MergeStopSequences(client []string, preset []string) []string {
	seen := make(map[string]struct{}, len(client)+len(preset))
	out := make([]string, 0, len(client)+len(preset))
	for _, s := range client {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	for _, s := range preset {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	return out
}

func fallbackPlainChatPrompt(messages []*domain.AIChatMessage, genParams *domain.GenerationParams) string {
	var b strings.Builder
	for _, m := range messages {
		if m == nil {
			continue
		}

		var role string
		switch m.Role {
		case domain.AIChatMessageRoleSystem:
			role = "System"
		case domain.AIChatMessageRoleAssistant:
			role = "Assistant"
		default:
			role = "User"
		}

		b.WriteString(role)
		b.WriteString(": ")
		b.WriteString(m.Content)
		b.WriteString("\n")
	}

	if genParams != nil && len(genParams.Tools) > 0 {
		b.WriteString(fallbackToolsBlock(genParams.Tools))
	}
	b.WriteString("Assistant: ")

	return b.String()
}

func fallbackToolsBlock(tools []domain.Tool) string {
	var b strings.Builder
	b.WriteString("\nTools:\n")
	for _, t := range tools {
		b.WriteString("- ")
		b.WriteString(t.Name)
		if t.Description != "" {
			b.WriteString(": ")
			b.WriteString(t.Description)
		}

		if t.ParametersJSON != "" {
			b.WriteString(" (params: ")
			b.WriteString(t.ParametersJSON)
			b.WriteString(")")
		}
		b.WriteString("\n")
	}
	b.WriteString("\nReply with JSON: {\"name\": \"tool_name\", \"arguments\": {...}}\n\n")

	return b.String()
}
