package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24server"
	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8786", "адрес HTTP (POST JSON-RPC, при необходимости GET SSE)")
	flag.Parse()
	webhookBase := os.Getenv("B24_WEBHOOK_BASE")
	logLevel := getenvDefault("B24_LOG_LEVEL", "info")
	retryMax := getenvIntDefault("B24_RETRY_MAX", 1)
	retryBackoffMS := getenvIntDefault("B24_RETRY_BACKOFF_MS", 300)
	log.Printf("MCP Bitrix24 streamable: starting webhook_base_set=%t listen=%s", webhookBase != "", *addr)

	srv, err := bitrix24server.NewServer(bitrix24server.Config{
		WebhookBase:    webhookBase,
		LogLevel:       logLevel,
		RetryMax:       retryMax,
		RetryBackoffMS: retryBackoffMS,
	})
	if err != nil {
		log.Fatalf("init bitrix24 server: %v", err)
	}

	h := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return srv
	}, nil)

	log.Printf("MCP Bitrix24 streamable: transport=streamable url=http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, mcpsafe.RecoverPanic("mcp-bitrix24-streamable", h)))
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
