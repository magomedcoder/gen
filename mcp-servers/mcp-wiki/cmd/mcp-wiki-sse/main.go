package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/magomedcoder/gen/mcp-servers/mcp-wiki/internal/mcpwikiserver"
	"github.com/magomedcoder/gen/pkg/mcpcache"
	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:8771", "адрес HTTP (GET = SSE, POST = сообщения сессии)")
	wikiDir := flag.String("wiki-dir", "", "обязательный каталог wiki для index_wiki_folder (единственный корень индексации)")
	flag.Parse()
	if strings.TrimSpace(*wikiDir) == "" {
		log.Fatal("mcp-wiki-sse: обязателен флаг -wiki-dir (каталог wiki)")
	}

	cache := mcpcache.NewServerByKey(func(key string) *mcp.Server {
		return mcpwikiserver.NewServerWithOptions(mcpwikiserver.Options{
			DefaultDirectory: key,
		})
	})

	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return cache.Get(*wikiDir)
	}, nil)

	log.Printf("MCP wiki server (sse): transport=sse url=http://%s/ default_wiki_dir=%q", *addr, strings.TrimSpace(*wikiDir))
	log.Fatal(http.ListenAndServe(*addr, mcpsafe.RecoverPanic("mcp-wiki-sse", handler)))
}
