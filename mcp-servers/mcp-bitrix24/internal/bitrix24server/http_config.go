package bitrix24server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	headerWebhookBase           = "X-B24-Base"
	headerLogLevel              = "X-B24-Log-Level"
	headerRetryMax              = "X-B24-Retry-Max"
	headerRetryBackoffMS        = "X-B24-Retry-Backoff-Ms"
	headerDisableHeavyAnalytics = "X-B24-Disable-Heavy-Analytics"
)

func ConfigFromHTTPRequest(r *http.Request, defaults Config) (Config, error) {
	cfg := defaults
	if r == nil {
		return cfg, nil
	}

	if v := pickRequestValue(r, headerWebhookBase, "b24_webhook_base"); v != "" {
		cfg.WebhookBase = v
	}

	if v := pickRequestValue(r, headerLogLevel, "b24_log_level"); v != "" {
		cfg.LogLevel = v
	}

	if v := pickRequestValue(r, headerRetryMax, "b24_retry_max"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("%s: expected integer, got %q", headerRetryMax, v)
		}
		cfg.RetryMax = n
	}

	if v := pickRequestValue(r, headerRetryBackoffMS, "b24_retry_backoff_ms"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("%s: expected integer, got %q", headerRetryBackoffMS, v)
		}
		cfg.RetryBackoffMS = n
	}

	if v := pickRequestValue(r, headerDisableHeavyAnalytics, "b24_disable_heavy_analytics"); v != "" {
		b, err := strconv.ParseBool(strings.ToLower(v))
		if err != nil {
			return cfg, fmt.Errorf("%s: expected bool, got %q", headerDisableHeavyAnalytics, v)
		}
		cfg.DisableHeavyAnalytics = b
	}

	return cfg, nil
}

func ConfigCacheKey(cfg Config) string {
	return strings.Join([]string{
		strings.TrimSpace(cfg.WebhookBase),
		strings.ToLower(strings.TrimSpace(cfg.LogLevel)),
		strconv.Itoa(cfg.RetryMax),
		strconv.Itoa(cfg.RetryBackoffMS),
		strconv.FormatBool(cfg.DisableHeavyAnalytics),
	}, "|")
}

func pickRequestValue(r *http.Request, headerName, queryName string) string {
	if r == nil {
		return ""
	}

	if v := strings.TrimSpace(r.Header.Get(headerName)); v != "" {
		return v
	}

	if queryName == "" || r.URL == nil {
		return ""
	}

	return strings.TrimSpace(r.URL.Query().Get(queryName))
}
