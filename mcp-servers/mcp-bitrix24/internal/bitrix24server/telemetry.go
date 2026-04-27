package bitrix24server

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	toolCallsTotal     uint64
	toolErrorsTotal    uint64
	toolLatencyUSum    uint64
	httpCallsTotal     uint64
	httpErrorsTotal    uint64
	softSkipsTotal     uint64
	telemetryStartOnce sync.Once
)

func initTelemetryReporter() {
	telemetryStartOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				log.Printf(
					"[b24-mcp] metrics tool_calls=%d tool_errors=%d avg_tool_latency_ms=%.2f http_calls=%d http_errors=%d soft_skips=%d",
					atomic.LoadUint64(&toolCallsTotal),
					atomic.LoadUint64(&toolErrorsTotal),
					avgToolLatencyMS(),
					atomic.LoadUint64(&httpCallsTotal),
					atomic.LoadUint64(&httpErrorsTotal),
					atomic.LoadUint64(&softSkipsTotal),
				)
			}
		}()
	})
}

func avgToolLatencyMS() float64 {
	calls := atomic.LoadUint64(&toolCallsTotal)
	if calls == 0 {
		return 0
	}

	return float64(atomic.LoadUint64(&toolLatencyUSum)) / float64(calls) / 1000.0
}

func incHTTPCall() {
	atomic.AddUint64(&httpCallsTotal, 1)
}

func incHTTPError() {
	atomic.AddUint64(&httpErrorsTotal, 1)
}

func incSoftSkip() {
	atomic.AddUint64(&softSkipsTotal, 1)
}

func incToolCall(d time.Duration) {
	atomic.AddUint64(&toolCallsTotal, 1)
	atomic.AddUint64(&toolLatencyUSum, uint64(d.Microseconds()))
}

func incToolError() {
	atomic.AddUint64(&toolErrorsTotal, 1)
}
