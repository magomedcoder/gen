package service

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/magomedcoder/gen/llm-runner/domain"
)

func messageHasPayload(m *domain.AIChatMessage) bool {
	if m == nil {
		return false
	}

	if len(m.AttachmentContent) > 0 {
		return true
	}

	if strings.TrimSpace(m.Content) != "" {
		return true
	}

	if m.Role == domain.AIChatMessageRoleAssistant && strings.TrimSpace(m.ToolCallsJSON) != "" {
		return true
	}

	if m.Role == domain.AIChatMessageRoleTool && strings.TrimSpace(m.ToolCallID) != "" {
		return true
	}

	return false
}

func FormatContentForBuiltinChatTemplate(m *domain.AIChatMessage) string {
	if m == nil {
		return ""
	}

	c := m.Content
	if m.Role == domain.AIChatMessageRoleTool {
		var b strings.Builder
		if m.ToolCallID != "" {
			fmt.Fprintf(&b, "[call_id=%s] ", m.ToolCallID)
		}

		if m.ToolName != "" {
			fmt.Fprintf(&b, "[%s] ", m.ToolName)
		}

		b.WriteString(c)
		return b.String()
	}

	if m.Role == domain.AIChatMessageRoleAssistant && strings.TrimSpace(m.ToolCallsJSON) != "" {
		if strings.TrimSpace(c) != "" {
			return c + "\n[tool_calls]: " + m.ToolCallsJSON
		}

		return "[tool_calls]: " + m.ToolCallsJSON
	}

	return c
}

func isImageFilenameExt(name string) bool {
	switch strings.ToLower(filepath.Ext(strings.TrimSpace(name))) {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func messageLikelyVisionImageAttachment(m *domain.AIChatMessage) bool {
	if m == nil || len(m.AttachmentContent) == 0 {
		return false
	}

	mt := strings.ToLower(strings.TrimSpace(m.AttachmentMime))
	if strings.HasPrefix(mt, "image/") {
		return true
	}

	if mt == "" {
		return isImageFilenameExt(m.AttachmentName) || strings.TrimSpace(m.AttachmentName) == ""
	}

	return false
}

func MessagesHaveVisionAttachments(messages []*domain.AIChatMessage) bool {
	return slices.ContainsFunc(messages, messageLikelyVisionImageAttachment)
}

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

		if !messageHasPayload(m) {
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
