package bitrix24server

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "b24_request_id"

var requestSeq uint64

func withRequestID(ctx context.Context) (context.Context, string) {
	if existing, ok := requestIDFromContext(ctx); ok && existing != "" {
		return ctx, existing
	}

	seq := atomic.AddUint64(&requestSeq, 1)
	rid := fmt.Sprintf("b24-%d-%d", time.Now().UnixMilli(), seq)
	return context.WithValue(ctx, requestIDKey, rid), rid
}

func requestIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v := ctx.Value(requestIDKey)
	s, ok := v.(string)
	return s, ok
}
