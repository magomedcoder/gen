package service

import (
	"strings"

	"github.com/magomedcoder/llm-runner/domain"
)

func buildPromptForTestFormat(format string, messages []*domain.AIChatMessage, genParams *domain.GenerationParams) string {
	if len(messages) == 0 {
		return ""
	}
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "llama3":
		return buildPromptLlama3Test(messages, genParams)
	default:
		return fallbackPlainChatPrompt(messages, genParams)
	}
}

const (
	testLlama3Begin = "<|begin_of_text|>"
	testLlama3HdrS  = "<|start_header_id|>"
	testLlama3HdrE  = "<|end_header_id|>"
	testLlama3EOT   = "<|eot_id|>"
)

func buildPromptLlama3Test(messages []*domain.AIChatMessage, genParams *domain.GenerationParams) string {
	var b strings.Builder
	b.WriteString(testLlama3Begin)
	for _, m := range messages {
		if m == nil {
			continue
		}
		name := ChatRoleString(m.Role)
		b.WriteString(testLlama3HdrS)
		b.WriteString(name)
		b.WriteString(testLlama3HdrE)
		b.WriteString("\n\n")
		b.WriteString(m.Content)
		b.WriteString(testLlama3EOT)
	}
	if genParams != nil && len(genParams.Tools) > 0 {
		b.WriteString(fallbackToolsBlock(genParams.Tools))
	}

	b.WriteString(testLlama3HdrS)
	b.WriteString("assistant")
	b.WriteString(testLlama3HdrE)
	b.WriteString("\n\n")

	return b.String()
}
