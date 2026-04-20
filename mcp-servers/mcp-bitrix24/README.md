### MCP Bitrix24 - transport stdio

Сервер для работы с задачами Bitrix24:
- список задач (`tasks.task.list`)
- получение задачи (`tasks.task.get`)
- получение комментариев задачи (`tasks.task.commentitem.getlist`)
- базовый анализ задачи по дедлайну, статусу и активности

Переменные окружения:
- `B24_WEBHOOK_BASE` - базовый webhook URL, например `https://bitrix24.example.com/rest/43176/00000000`

Сборка:

```bash
go build -o mcp-bitrix24-stdio ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-stdio
```

```
transport = stdio

command = путь к бинарнику

args` пустые
```

---

### MCP Bitrix24 - transport SSE

Сборка:

```bash
go build -o mcp-bitrix24-sse ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-sse
```

Запуск:

```bash
B24_WEBHOOK_BASE="https://bitrix24.example.com/rest/43176/00000000" ./mcp-bitrix24-sse -listen 127.0.0.1:8785
```

`transport = sse`, `url = http://127.0.0.1:8785/`

---

### MCP Bitrix24 - transport streamable HTTP

Сборка:

```bash
go build -o mcp-bitrix24-streamable ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-streamable
```

Запуск:

```bash
B24_WEBHOOK_BASE="https://bitrix24.example.com/rest/43176/00000000" ./mcp-bitrix24-streamable -listen 127.0.0.1:8786
```

`transport = streamable`, `url = http://127.0.0.1:8786/`

---

### Инструменты MCP

Краткий список (схемы аргументов задаёт клиент MCP по полям инструментов):

- `b24_list_tasks` - `tasks.task.list` (`filter`, `select`, `order`, `start`)
- `b24_get_task` - `tasks.task.get` (`task_id`, `select`)
- `b24_get_task_comments` - `tasks.task.commentitem.getlist` (`task_id`, `order`, `select`)
- `b24_analyze_task` - сводка по задаче и комментариям (`task_id`, `include_comments`)
- `b24_call_method` - произвольный REST-метод (`method`, `params`)

---

```bash
go build -o mcp-bitrix24-mock-rest ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-mock-rest

./mcp-bitrix24-mock-rest -listen 127.0.0.1:8899
```
