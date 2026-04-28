package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/magomedcoder/gen/mcp-servers/mcp-bitrix24/internal/bitrix24server"
	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8785", "адрес HTTP (GET = SSE, POST = сообщения сессии)")
	flag.Parse()
	defaultCfg := bitrix24server.Config{}
	log.Printf("MCP Bitrix24 SSE: starting listen=%s", *addr)

	var (
		mu      sync.Mutex
		servers = map[string]*mcp.Server{}
	)

	h := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		cfg, err := bitrix24server.ConfigFromHTTPRequest(r, defaultCfg)
		if err != nil {
			log.Printf("MCP Bitrix24 SSE: invalid request config: %v; fallback to defaults", err)
			cfg = defaultCfg
		}

		key := bitrix24server.ConfigCacheKey(cfg)
		mu.Lock()
		defer mu.Unlock()

		if srv, ok := servers[key]; ok {
			return srv
		}

		srv, err := bitrix24server.NewServer(cfg)
		if err != nil {
			log.Printf("MCP Bitrix24 SSE: init server by request config failed: %v; fallback to defaults", err)
			srv, err = bitrix24server.NewServer(defaultCfg)
			if err != nil {
				log.Printf("MCP Bitrix24 SSE: fallback init failed: %v", err)
				return mcp.NewServer(&mcp.Implementation{Name: "bitrix24", Version: "1.0.0"}, nil)
			}
		}

		servers[key] = srv
		return srv
	}, nil)

	log.Printf("MCP Bitrix24 SSE: transport=sse url=http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, mcpsafe.RecoverPanic("mcp-bitrix24-sse", h)))
}
