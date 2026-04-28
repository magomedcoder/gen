package bitrix24server

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24mock"
)

func TestBuildAnalyticsContextForTaskList_WithMock(t *testing.T) {
	mock := bitrix24mock.NewServer()
	srv := httptest.NewServer(mock.Handler())
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, 2*time.Second, "info", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	ac, err := buildAnalyticsContextForTaskList(context.Background(), client, nil, nil, nil, 10, false)
	if err != nil {
		t.Fatalf("buildAnalyticsContextForTaskList: %v", err)
	}

	if ac == nil {
		t.Fatalf("expected non-nil analytics context")
	}

	if len(ac.Items) == 0 {
		t.Fatalf("expected non-empty items from mock")
	}
}
