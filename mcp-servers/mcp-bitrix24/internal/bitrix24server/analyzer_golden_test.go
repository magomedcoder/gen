package bitrix24server

import (
	"strings"
	"testing"
	"time"
)

func TestTaskConclusion_GoldenOverdueBlocked(t *testing.T) {
	now := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	deadline := now.Add(-2 * time.Hour)

	got := taskConclusion("Инцидент 42", 3, "высокий", deadline, now, 5, true)

	wantContains := []string{
		"Статус: критично",
		"Главный риск: задача \"Инцидент 42\" просрочена и заблокирована",
		"Приоритетное действие: срочная эскалация, снять блокер и зафиксировать новый план",
		"Срок реакции: немедленно",
	}

	for _, s := range wantContains {
		if !strings.Contains(got, s) {
			t.Fatalf("expected conclusion to contain %q, got:\n%s", s, got)
		}
	}
}

func TestTaskConclusion_GoldenStable(t *testing.T) {
	now := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	deadline := now.Add(24 * time.Hour)

	got := taskConclusion("План релиза", 3, "низкий", deadline, now, 3, false)
	if !strings.Contains(got, "Статус: стабильно") {
		t.Fatalf("expected stable status, got:\n%s", got)
	}

	if !strings.Contains(got, "Срок реакции: по текущему регламенту") {
		t.Fatalf("expected reaction SLA line, got:\n%s", got)
	}
}
