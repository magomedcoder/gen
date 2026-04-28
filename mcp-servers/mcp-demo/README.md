# MCP demo

Минимальный пример MCP-сервера: SSE и streamable HTTP.

---

## Transport SSE

Сборка:

```bash
go build -o ./build/mcp-sse-demo ./mcp-servers/mcp-demo/cmd/mcp-sse-demo
```

Запуск:

```bash
./build/mcp-sse-demo -listen 127.0.0.1:8765
```

В GEN: `transport = sse`, `url = http://127.0.0.1:8765/`

Для удалённого хоста добавьте его в mcp.http_allow_hosts (или http_allow_any)

---

## Transport streamable HTTP

Сборка:

```bash
go build -o ./build/mcp-streamable-demo ./mcp-servers/mcp-demo/cmd/mcp-streamable-demo
```

Запуск:

```bash
./build/mcp-streamable-demo -listen 127.0.0.1:8766
```

```
transport = streamable

url = http://127.0.0.1:8766/
```
