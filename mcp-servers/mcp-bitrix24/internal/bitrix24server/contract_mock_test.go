package bitrix24server

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24mock"
)

func TestContract_LoadTaskList_WithMockServer(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	tasks, err := loadTaskList(context.Background(), client, nil, nil, nil, 50)
	if err != nil {
		t.Fatalf("loadTaskList: %v", err)
	}

	if len(tasks) == 0 {
		t.Fatalf("expected non-empty tasks list from mock")
	}
}

func TestContract_LoadTaskComments_WithMockServer(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	comments, err := loadTaskComments(context.Background(), client, 1001)
	if err != nil {
		t.Fatalf("loadTaskComments: %v", err)
	}

	if len(comments) == 0 {
		t.Fatalf("expected non-empty comments for task 1001 from mock")
	}
}

func TestContract_RunExecutiveSummary_WithMockServer(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	report, err := runExecutiveSummary(context.Background(), client, nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("runExecutiveSummary: %v", err)
	}

	if !strings.Contains(report, "Executive summary") {
		t.Fatalf("unexpected report body: %s", report)
	}

	if !strings.Contains(report, "=== Вывод ===") {
		t.Fatalf("expected conclusion block in report: %s", report)
	}
}
