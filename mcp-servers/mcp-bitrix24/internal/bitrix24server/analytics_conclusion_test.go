package bitrix24server

import (
	"strings"
	"testing"
)

func TestFormatConclusion_GoldenShape(t *testing.T) {
	got := formatConclusion("критично", "просрочено 2", "эскалация", "сегодня")
	want := "Статус: критично\nГлавный риск: просрочено 2\nПриоритетное действие: эскалация\nСрок реакции: сегодня"
	if got != want {
		t.Fatalf("unexpected formatConclusion output\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}

func TestSLAConclusion_GoldenScenarios(t *testing.T) {
	cases := []struct {
		name     string
		overdue  int
		today    int
		soon     int
		noDL     int
		contains string
	}{
		{"overdue", 3, 0, 0, 0, "Статус: критично"},
		{"today", 0, 2, 0, 0, "Статус: на грани"},
		{"no-deadline", 0, 0, 0, 4, "Статус: под контролем"},
		{"green", 0, 0, 0, 0, "Статус: стабильно"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := slaConclusion(tc.overdue, tc.today, tc.soon, tc.noDL)
			if !strings.Contains(got, tc.contains) {
				t.Fatalf("expected %q in output, got:\n%s", tc.contains, got)
			}

			if !strings.Contains(got, "Срок реакции:") {
				t.Fatalf("expected reaction line in output, got:\n%s", got)
			}
		})
	}
}
