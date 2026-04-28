package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	webhookBase := os.Getenv("B24_WEBHOOK_BASE")
	logLevel := getenvDefault("B24_LOG_LEVEL", "info")
	retryMax := getenvIntDefault("B24_RETRY_MAX", 1)
	retryBackoffMS := getenvIntDefault("B24_RETRY_BACKOFF_MS", 300)
	disableHeavyAnalytics := getenvBoolDefault("B24_DISABLE_HEAVY_ANALYTICS", false)
	log.Printf("MCP Bitrix24 stdio: starting webhook_base_set=%t", webhookBase != "")
	srv, err := bitrix24server.NewServer(bitrix24server.Config{
		WebhookBase:           webhookBase,
		LogLevel:              logLevel,
		RetryMax:              retryMax,
		RetryBackoffMS:        retryBackoffMS,
		DisableHeavyAnalytics: disableHeavyAnalytics,
	})
	if err != nil {
		log.Fatalf("init bitrix24 server: %v", err)
	}

	log.Printf("MCP Bitrix24 stdio: transport=stdio ready")
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("mcp server: %v", err)
	}
}

func getenvDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getenvIntDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getenvBoolDefault(key string, fallback bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return fallback
	}

	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
