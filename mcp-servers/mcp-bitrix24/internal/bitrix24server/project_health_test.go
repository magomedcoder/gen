package bitrix24server

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24mock"
)

func TestRunProjectHealth_WithMockServer(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	report, err := runProjectHealth(context.Background(), client, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("runProjectHealth: %v", err)
	}

	if !strings.Contains(report, "Project health summary") {
		t.Fatalf("unexpected report: %s", report)
	}

	if !strings.Contains(report, "Health score:") {
		t.Fatalf("missing health score: %s", report)
	}

	if !strings.Contains(report, "=== Вывод ===") {
		t.Fatalf("missing conclusion: %s", report)
	}
}
