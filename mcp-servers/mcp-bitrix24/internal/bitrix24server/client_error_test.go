package bitrix24server

import (
	"strings"
	"testing"
)

func TestMaskJSONForLog_MasksSensitiveFields(t *testing.T) {
	raw := []byte(`{"auth":"abc","nested":{"token":"secret","ok":"value"},"password":"123"}`)
	masked := string(maskJSONForLog(raw))

	if strings.Contains(masked, `"abc"`) || strings.Contains(masked, `"secret"`) || strings.Contains(masked, `"123"`) {
		t.Fatalf("sensitive values were not masked: %s", masked)
	}

	if !strings.Contains(masked, "***MASKED***") {
		t.Fatalf("expected masked marker in output: %s", masked)
	}
}

func TestBitrixErrorHint_TaskCommentActionFailed(t *testing.T) {
	hint := bitrixErrorHint("task.commentitem.getlist", "ERROR_CORE", "TASKS_ERROR_EXCEPTION_#8; Action failed", 400)
	if !strings.Contains(strings.ToLower(hint), "include_comments") {
		t.Fatalf("expected include_comments hint, got: %s", hint)
	}
}

func TestWrapBitrixError_ContainsMethodAndHint(t *testing.T) {
	raw := []byte(`{"error":"ERROR_CORE","error_description":"TASKS_ERROR_EXCEPTION_#256; expected to be of type integer"}`)
	err := wrapBitrixError("task.commentitem.getlist", 400, raw, nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "task.commentitem.getlist") {
		t.Fatalf("expected method name in error, got: %s", msg)
	}

	if !strings.Contains(strings.ToLower(msg), "проверьте типы") {
		t.Fatalf("expected actionable hint in error, got: %s", msg)
	}
}
