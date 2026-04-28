package bitrix24server

import "testing"

func TestExtractTaskListPage_SupportsDifferentShapes(t *testing.T) {
	respWithTasks := map[string]any{
		"result": map[string]any{
			"tasks": []any{
				map[string]any{"ID": "1001", "TITLE": "A"},
			},
		},
		"next": float64(50),
	}

	tasks, next := extractTaskListPage(respWithTasks)
	if len(tasks) != 1 || numberLike(tasks[0]["ID"]) != 1001 {
		t.Fatalf("unexpected tasks from result.tasks: %+v", tasks)
	}

	if next == nil || *next != 50 {
		t.Fatalf("unexpected next from result.tasks: %v", next)
	}

	respWithArrayResult := map[string]any{
		"result": []any{
			map[string]any{"ID": "2001", "TITLE": "B"},
		},
	}

	tasks, next = extractTaskListPage(respWithArrayResult)
	if len(tasks) != 1 || numberLike(tasks[0]["ID"]) != 2001 {
		t.Fatalf("unexpected tasks from result array: %+v", tasks)
	}

	if next != nil {
		t.Fatalf("expected nil next for array result, got %v", *next)
	}
}

func TestExtractNextStart_ParsesSupportedTypes(t *testing.T) {
	v := extractNextStart(map[string]any{"next": float64(120)})
	if v == nil || *v != 120 {
		t.Fatalf("float64 next parse failed: %v", v)
	}

	v = extractNextStart(map[string]any{"next": "75"})
	if v == nil || *v != 75 {
		t.Fatalf("string next parse failed: %v", v)
	}

	v = extractNextStart(map[string]any{"next": "bad"})
	if v != nil {
		t.Fatalf("invalid string next must be nil, got %v", *v)
	}
}
