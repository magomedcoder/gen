package bitrix24server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func TestBitrixClientCall_DecodesWindows1251JSON(t *testing.T) {
	cp1251Body, err := charmap.Windows1251.NewEncoder().Bytes([]byte(`{"result":{"task":{"TITLE":"Привет"}}}`))
	if err != nil {
		t.Fatalf("encode cp1251: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=windows-1251")
		_, _ = w.Write(cp1251Body)
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second, "debug", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	resp, err := client.call(context.Background(), "tasks.task.get", map[string]any{"taskId": 1001})
	if err != nil {
		t.Fatalf("call: %v", err)
	}

	got := nestedTaskTitle(resp)
	if got != "Привет" {
		t.Fatalf("unexpected title: %q", got)
	}
}

func TestBitrixClientCall_DecodesWindows1251WithoutCharsetHint(t *testing.T) {
	cp1251Body, err := charmap.Windows1251.NewEncoder().Bytes([]byte(`{"result":{"task":{"TITLE":"Задача"}}}`))
	if err != nil {
		t.Fatalf("encode cp1251: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(cp1251Body)
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second, "debug", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	resp, err := client.call(context.Background(), "tasks.task.get", map[string]any{"taskId": 1001})
	if err != nil {
		t.Fatalf("call: %v", err)
	}

	got := nestedTaskTitle(resp)
	if got != "Задача" {
		t.Fatalf("unexpected title: %q", got)
	}
}

func nestedTaskTitle(resp map[string]any) string {
	result, _ := resp["result"].(map[string]any)
	task, _ := result["task"].(map[string]any)
	title, _ := task["TITLE"].(string)
	return title
}

func TestCallTaskCommentItemGetList_UsesStablePayloadOrder(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		gotBody = string(raw)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":[]}`))
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second, "debug", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	_, err = client.callTaskCommentItemGetList(
		context.Background(),
		1822404,
		map[string]any{"POST_DATE": "desc"},
		map[string]any{"AUTHOR_ID": 7},
	)
	if err != nil {
		t.Fatalf("callTaskCommentItemGetList: %v", err)
	}

	idxTask := strings.Index(gotBody, `"TASKID"`)
	idxOrder := strings.Index(gotBody, `"ORDER"`)
	idxFilter := strings.Index(gotBody, `"FILTER"`)
	if idxTask == -1 || idxOrder == -1 || idxFilter == -1 {
		t.Fatalf("payload missing required keys: %s", gotBody)
	}
	if !(idxTask < idxOrder && idxOrder < idxFilter) {
		t.Fatalf("unexpected key order, want TASKID->ORDER->FILTER, got: %s", gotBody)
	}
}

func TestBitrixClientCall_BlocksWriteMethodsByReadOnlyPolicy(t *testing.T) {
	var hitCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hitCount, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":{"ok":true}}`))
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second, "debug", 0, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	_, err = client.call(context.Background(), "tasks.task.update", map[string]any{
		"taskId": 1001,
	})
	if err == nil {
		t.Fatalf("expected read-only policy error")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "read-only policy") {
		t.Fatalf("unexpected error text: %v", err)
	}

	if atomic.LoadInt32(&hitCount) != 0 {
		t.Fatalf("write method must be blocked before HTTP call, hit_count=%d", hitCount)
	}
}
