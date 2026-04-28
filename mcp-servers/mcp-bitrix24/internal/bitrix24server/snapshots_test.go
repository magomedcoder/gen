package bitrix24server

import "testing"

func TestNormalizeTaskSnapshot(t *testing.T) {
	task := map[string]any{
		"ID":                 "1001",
		"TITLE":              "Demo task",
		"STATUS":             "3",
		"CREATED_DATE":       "2026-04-28T10:00:00+03:00",
		"RESPONSIBLE_ID":     "21",
		"TIME_ESTIMATE":      "3600",
		"TIME_SPENT_IN_LOGS": "1200",
	}

	s := normalizeTaskSnapshot(task)
	if s.ID != 1001 {
		t.Fatalf("unexpected id: %d", s.ID)
	}

	if s.Title != "Demo task" {
		t.Fatalf("unexpected title: %q", s.Title)
	}

	if s.StatusCode != 3 || s.StatusLabel == "" {
		t.Fatalf("unexpected status: code=%d label=%q", s.StatusCode, s.StatusLabel)
	}

	if s.ResponsibleID != "21" {
		t.Fatalf("unexpected responsible: %q", s.ResponsibleID)
	}
}

func TestNormalizeCommentSnapshots_UsesFallbackTaskID(t *testing.T) {
	comments := []map[string]any{
		{
			"ID":           "9001",
			"AUTHOR_ID":    "17",
			"POST_DATE":    "2026-04-28T11:00:00+03:00",
			"POST_MESSAGE": "hello",
		},
	}

	out := normalizeCommentSnapshots(comments, 1001)
	if len(out) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(out))
	}

	if out[0].TaskID != 1001 {
		t.Fatalf("unexpected fallback task id: %d", out[0].TaskID)
	}

	if out[0].Message != "hello" {
		t.Fatalf("unexpected message: %q", out[0].Message)
	}
}

func TestBuildTaskTimeline_SortsByDateDesc(t *testing.T) {
	task := TaskSnapshot{
		ID:          1001,
		Title:       "Demo",
		StatusCode:  3,
		StatusLabel: "Выполняется",
		CreatedAt:   "2026-04-28T10:00:00+03:00",
		ChangedAt:   "2026-04-28T12:00:00+03:00",
	}
	comments := []CommentSnapshot{
		{
			ID:        1,
			TaskID:    1001,
			CreatedAt: "2026-04-28T13:00:00+03:00",
			Message:   "latest",
		},
		{
			ID:        2,
			TaskID:    1001,
			CreatedAt: "2026-04-28T11:00:00+03:00",
			Message:   "old",
		},
	}

	events := buildTaskTimeline(task, comments)
	if len(events) < 3 {
		t.Fatalf("expected at least 3 events, got %d", len(events))
	}

	if events[0].Type != "comment_added" || events[0].Details != "latest" {
		t.Fatalf("unexpected first event: %+v", events[0])
	}
}
