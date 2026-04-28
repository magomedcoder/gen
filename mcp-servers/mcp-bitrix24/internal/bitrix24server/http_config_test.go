package bitrix24server

import (
	"net/http/httptest"
	"testing"
)

func TestConfigFromHTTPRequest_HeadersOverrideDefaults(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1:8785/?b24_retry_max=3", nil)
	req.Header.Set(headerWebhookBase, "https://example.bitrix24.ru/rest/1/token")
	req.Header.Set(headerRetryMax, "7")
	req.Header.Set(headerRetryBackoffMS, "900")
	req.Header.Set(headerDisableHeavyAnalytics, "true")
	req.Header.Set(headerLogLevel, "debug")

	cfg, err := ConfigFromHTTPRequest(req, Config{
		WebhookBase:           "https://default/rest/1/default",
		LogLevel:              "info",
		RetryMax:              1,
		RetryBackoffMS:        300,
		DisableHeavyAnalytics: false,
	})
	if err != nil {
		t.Fatalf("ConfigFromHTTPrequest: %v", err)
	}

	if cfg.WebhookBase != "https://example.bitrix24.ru/rest/1/token" {
		t.Fatalf("неожиданный webhook_base: %q", cfg.WebhookBase)
	}
	if cfg.LogLevel != "debug" || cfg.RetryMax != 7 || cfg.RetryBackoffMS != 900 || !cfg.DisableHeavyAnalytics {
		t.Fatalf("неожиданный cfg: %+v", cfg)
	}
}

func TestConfigFromHTTPRequest_QueryFallback(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1:8785/?b24_webhook_base=https://tenant.bitrix24.ru/rest/2/token&b24_retry_max=4&b24_retry_backoff_ms=500&b24_disable_heavy_analytics=true&b24_log_level=debug", nil)

	cfg, err := ConfigFromHTTPRequest(req, Config{
		LogLevel:       "info",
		RetryMax:       1,
		RetryBackoffMS: 300,
	})
	if err != nil {
		t.Fatalf("ConfigFromHTTPrequest: %v", err)
	}

	if cfg.WebhookBase != "https://tenant.bitrix24.ru/rest/2/token" {
		t.Fatalf("неожиданный webhook_base: %q", cfg.WebhookBase)
	}
	if cfg.LogLevel != "debug" || cfg.RetryMax != 4 || cfg.RetryBackoffMS != 500 || !cfg.DisableHeavyAnalytics {
		t.Fatalf("неожиданный cfg: %+v", cfg)
	}
}

func TestConfigFromHTTPRequest_InvalidValue(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1:8785/", nil)
	req.Header.Set(headerRetryMax, "oops")

	_, err := ConfigFromHTTPRequest(req, Config{})
	if err == nil {
		t.Fatal("expected parse error for invalid retry max")
	}
}
