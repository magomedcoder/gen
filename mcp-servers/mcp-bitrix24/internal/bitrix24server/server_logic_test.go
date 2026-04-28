package bitrix24server

import "testing"

func TestExtractTaskIDFromListFilter(t *testing.T) {
	taskID, ok := extractTaskIDFromListFilter(map[string]any{"ID": "100"})
	if !ok || taskID != 100 {
		t.Fatalf("expected task id 100, got ok=%v id=%d", ok, taskID)
	}

	taskID, ok = extractTaskIDFromListFilter(map[string]any{"=ID": 42})
	if !ok || taskID != 42 {
		t.Fatalf("expected task id 42 from =ID, got ok=%v id=%d", ok, taskID)
	}

	taskID, ok = extractTaskIDFromListFilter(map[string]any{"taskId": " 777 "})
	if !ok || taskID != 777 {
		t.Fatalf("expected task id 777 from taskId, got ok=%v id=%d", ok, taskID)
	}

	if _, ok := extractTaskIDFromListFilter(map[string]any{"ID": "abc"}); ok {
		t.Fatalf("did not expect task-id detection for invalid ID value")
	}

	if _, ok := extractTaskIDFromListFilter(map[string]any{"RESPONSIBLE_ID": 7}); ok {
		t.Fatalf("did not expect task-id detection for non-id filter")
	}
}

func TestExtractComments_FromMapPayloadShapes(t *testing.T) {
	respWithItemsMap := map[string]any{
		"result": map[string]any{
			"items": map[string]any{
				"9001": map[string]any{
					"ID":           "9001",
					"AUTHOR_ID":    "17",
					"POST_MESSAGE": "msg1",
				},
			},
		},
	}
	comments := extractComments(respWithItemsMap)
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment from items map, got %d", len(comments))
	}

	respWithCommentsObject := map[string]any{
		"result": map[string]any{
			"comments": map[string]any{
				"ID":           "9002",
				"AUTHOR_ID":    "21",
				"POST_MESSAGE": "msg2",
			},
		},
	}
	comments = extractComments(respWithCommentsObject)
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment from single comment object, got %d", len(comments))
	}

	respWithCommentsMapByID := map[string]any{
		"result": map[string]any{
			"comments": map[string]any{
				"9003": map[string]any{
					"ID":           "9003",
					"AUTHOR_ID":    "5",
					"POST_MESSAGE": "msg3",
				},
				"9004": map[string]any{
					"ID":           "9004",
					"AUTHOR_ID":    "8",
					"POST_MESSAGE": "msg4",
				},
			},
		},
	}
	comments = extractComments(respWithCommentsMapByID)
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments from comments map by id, got %d", len(comments))
	}
}

func TestExtractComments_FromNestedArraysInsideMap(t *testing.T) {
	resp := map[string]any{
		"result": map[string]any{
			"comments": map[string]any{
				"group_a": []any{
					map[string]any{"ID": "1", "AUTHOR_ID": "7", "POST_MESSAGE": "m1"},
				},
				"group_b": []any{
					map[string]any{"ID": "2", "AUTHOR_ID": "8", "POST_MESSAGE": "m2"},
				},
			},
		},
	}

	comments := extractComments(resp)
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments from nested arrays in map, got %d", len(comments))
	}
}

func TestCommentDiagnostics(t *testing.T) {
	resp := map[string]any{
		"result": map[string]any{
			"comments": map[string]any{
				"1": map[string]any{"ID": "1", "AUTHOR_ID": "7", "POST_MESSAGE": "ok"},
			},
		},
	}

	total := commentsTotalFromResponse(resp)
	if total != 1 {
		t.Fatalf("expected total=1, got %d", total)
	}

	warns := commentParseWarnings(2, 1)
	if len(warns) == 0 {
		t.Fatalf("expected warnings for partial parse")
	}
}
