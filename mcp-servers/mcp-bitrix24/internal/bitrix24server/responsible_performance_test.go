package bitrix24server

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24mock"
)

func TestRunResponsiblePerformance_WithMockServer(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	report, err := runResponsiblePerformance(context.Background(), client, "21", nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("runResponsiblePerformance: %v", err)
	}

	if !strings.Contains(report, "Performance по ответственному 21") {
		t.Fatalf("unexpected report: %s", report)
	}

	if !strings.Contains(report, "=== Вывод ===") {
		t.Fatalf("missing conclusion block: %s", report)
	}
}

func TestRunResponsiblePerformance_DoesNotMutateInputFilter(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	inputFilter := map[string]any{
		"STATUS": 3,
	}

	_, err = runResponsiblePerformance(context.Background(), client, "21", inputFilter, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("runResponsiblePerformance: %v", err)
	}

	if _, exists := inputFilter["RESPONSIBLE_ID"]; exists {
		t.Fatalf("input filter must not be mutated, got: %#v", inputFilter)
	}
}
