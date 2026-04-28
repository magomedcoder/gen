package bitrix24server

import (
	"testing"
	"time"
)

func TestBuildExecutionDriftReport_HighDriftOnOverrunAndSilence(t *testing.T) {
	now := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	task := TaskSnapshot{
		ID:           1001,
		TimeEstimate: 100,
		TimeSpent:    150,
	}
	comments := []CommentSnapshot{
		{
			CreatedAt: "2026-04-20T12:00:00Z",
			Message:   "old",
		},
	}

	report := buildExecutionDriftReport(task, comments, now)
	if report.DriftLevel != "high" {
		t.Fatalf("expected high drift, got %s", report.DriftLevel)
	}

	if report.OverrunSeconds <= 0 {
		t.Fatalf("expected overrun > 0, got %d", report.OverrunSeconds)
	}
}

func TestExecutionDriftActions_DefaultWhenLow(t *testing.T) {
	actions := executionDriftActions(ExecutionDriftReport{DriftLevel: "low"})
	if len(actions) == 0 {
		t.Fatalf("expected at least one action")
	}
}
