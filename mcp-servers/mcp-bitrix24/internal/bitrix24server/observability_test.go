package bitrix24server

import (
	"context"
	"testing"
)

func TestWithRequestID_AssignsAndReuses(t *testing.T) {
	ctx := context.Background()
	ctx1, id1 := withRequestID(ctx)
	if id1 == "" {
		t.Fatalf("expected non-empty request id")
	}

	got, ok := requestIDFromContext(ctx1)
	if !ok || got == "" {
		t.Fatalf("expected request id in context")
	}

	if got != id1 {
		t.Fatalf("unexpected id from context: %q != %q", got, id1)
	}

	ctx2, id2 := withRequestID(ctx1)
	if id2 != id1 {
		t.Fatalf("expected request id reuse, got %q want %q", id2, id1)
	}

	if ctx2 != ctx1 {
		t.Fatalf("expected same context when id already exists")
	}
}
