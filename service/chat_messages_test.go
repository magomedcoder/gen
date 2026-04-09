package service

import (
	"strings"
	"testing"

	"github.com/magomedcoder/gen-runner/domain"
)

func TestFormatContentForBuiltinChatTemplate_toolRoleLongMCPStyleContent(t *testing.T) {
	t.Parallel()
	const n = 200_000
	long := strings.Repeat("Z", n)
	m := &domain.AIChatMessage{
		Role:       domain.AIChatMessageRoleTool,
		Content:    long,
		ToolCallID: "call_abc",
		ToolName:   "mcp_3_h636174",
	}

	got := FormatContentForBuiltinChatTemplate(m)
	if !strings.HasPrefix(got, "[call_id=call_abc] ") {
		prefixLen := 40
		if len(got) < prefixLen {
			prefixLen = len(got)
		}
		t.Fatalf("prefix: %q", got[:prefixLen])
	}

	if !strings.Contains(got, "[mcp_3_h636174] ") {
		t.Fatal("ожидалось имя инструмента в префиксе")
	}

	if len(got) != len("[call_id=call_abc] [mcp_3_h636174] ")+n {
		t.Fatalf("длина: got %d want %d", len(got), len("[call_id=call_abc] [mcp_3_h636174] ")+n)
	}

	if !strings.HasSuffix(got, "ZZZ") {
		t.Fatal("хвост контента обрезан")
	}
}
