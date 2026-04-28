package bitrix24server

import (
	"testing"
	"time"
)

func TestDetectBlockerSignals_FindsKeywordsAndSortsByAge(t *testing.T) {
	now := time.Date(2026, 4, 28, 15, 0, 0, 0, time.UTC)
	comments := []CommentSnapshot{
		{
			ID:        1,
			TaskID:    1001,
			AuthorID:  "7",
			CreatedAt: "2026-04-28T14:00:00Z",
			Message:   "Жду ответ от клиента",
		},
		{
			ID:        2,
			TaskID:    1001,
			AuthorID:  "8",
			CreatedAt: "2026-04-26T10:00:00Z",
			Message:   "Blocked by dependency",
		},
		{
			ID:        3,
			TaskID:    1001,
			AuthorID:  "9",
			CreatedAt: "2026-04-28T13:00:00Z",
			Message:   "обычный апдейт",
		},
	}

	signals := detectBlockerSignals(comments, now)
	if len(signals) != 2 {
		t.Fatalf("expected 2 blocker signals, got %d", len(signals))
	}

	if signals[0].CommentID != 2 {
		t.Fatalf("expected oldest blocker first, got comment_id=%d", signals[0].CommentID)
	}

	if signals[1].Keyword == "" {
		t.Fatalf("expected detected keyword")
	}
}

func TestBlockerOwners_UniqueSorted(t *testing.T) {
	owners := blockerOwners([]BlockerSignal{
		{AuthorID: "8"},
		{AuthorID: "7"},
		{AuthorID: "8"},
	})

	if len(owners) != 2 || owners[0] != "7" || owners[1] != "8" {
		t.Fatalf("unexpected owners: %#v", owners)
	}
}
