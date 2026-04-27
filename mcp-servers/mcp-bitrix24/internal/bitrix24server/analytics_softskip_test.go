package bitrix24server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRunAnalyticsQuery_CommentErrorDoesNotFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/tasks.task.list"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"result": {
					"tasks": [
						{
							"ID":"101",
							"TITLE":"Тестовая задача",
							"STATUS":"3",
							"CHANGED_DATE":"2026-04-27T10:00:00+00:00"
						}
					]
				}
			}`))
		case strings.HasSuffix(r.URL.Path, "/task.commentitem.getlist"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"ERROR_CORE","error_description":"TASKS_ERROR_EXCEPTION_#8; Action failed; 8/TE/ACTION_FAILED_TO_BE_PROCESSED<br>"}`))
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"NOT_FOUND","error_description":"unexpected method"}`))
		}
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	includeComments := true
	limit := 20
	result, err := runAnalyticsQuery(
		context.Background(),
		client,
		"общий аналитический обзор",
		nil,
		nil,
		nil,
		nil,
		&limit,
		&includeComments,
	)

	if err != nil {
		t.Fatalf("runAnalyticsQuery returned error, expected soft-skip behavior: %v", err)
	}

	if !strings.Contains(result, "Найдено задач: 1") {
		t.Fatalf("unexpected analytics output: %s", result)
	}
}

func TestIsIgnorableCommentError(t *testing.T) {
	err := wrapBitrixError(
		"task.commentitem.getlist",
		http.StatusBadRequest,
		[]byte(`{"error":"ERROR_CORE","error_description":"TASKS_ERROR_EXCEPTION_#8; Action failed; 8/TE/ACTION_FAILED_TO_BE_PROCESSED<br>"}`),
		nil,
	)
	if !isIgnorableCommentError(err) {
		t.Fatalf("expected error to be ignorable for comments")
	}
}
