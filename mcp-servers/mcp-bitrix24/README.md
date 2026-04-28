# MCP Bitrix24

Сервер MCP для задач Bitrix24 через входящий webhook:

- список задач (`tasks.task.list`)
- получение задачи (`tasks.task.get`)
- комментарии (`task.commentitem.getlist`)
- таймлайн задачи (`b24_get_task_timeline`)
- анализ блокеров задачи (`b24_analyze_task_blockers`)
- анализ дрифта исполнения (`b24_analyze_task_execution_drift`)
- performance по ответственному (`b24_analyze_responsible_performance`)
- health-score проекта (`b24_analyze_project_health`)
- сводка по задаче (`b24_analyze_task`)
- аналитика по запросу (`b24_analyze_tasks_by_query`)
- портфельная аналитика (`b24_analyze_tasks_portfolio`)
- executive summary (`b24_analyze_tasks_executive_summary`)
- SLA-аналитика (`b24_analyze_tasks_sla`)

## Персональный конфиг через SSE/streamable

Поддерживаемые входные заголовки:

| Заголовок                       | Назначение                                                                       |
|---------------------------------|----------------------------------------------------------------------------------|
| `X-B24-Base`                    | Базовый URL webhook, например `https://bitrix24.example.com/rest/43176/00000000` |
| `X-B24-Log-Level`               | Уровень логов (`info`, `debug`)                                                  |
| `X-B24-Retry-Max`               | Число ретраев HTTP-запросов к Bitrix24                                           |
| `X-B24-Retry-Backoff-Ms`        | Базовый backoff в миллисекундах                                                  |
| `X-B24-Disable-Heavy-Analytics` | `true/false`, отключение тяжелой аналитики                                       |

Все параметры задаются через `headers` в конфиге MCP-сервера.

Пример записи MCP-сервера в GEN (SSE):

```json
{
  "transport": "sse",
  "url": "http://127.0.0.1:8785/",
  "headers": {
    "X-B24-Base": "https://bitrix24.example.com/rest/43176/00000000",
    "X-B24-Retry-Max": "2",
    "X-B24-Retry-Backoff-Ms": "400",
    "X-B24-Disable-Heavy-Analytics": "false"
  },
  "timeoutSeconds": 120
}
```

---

## Transport SSE

Сборка:

```bash
go build -o ./build/mcp-bitrix24-sse ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-sse
```

Запуск:

```bash
./build/mcp-bitrix24-sse -listen 127.0.0.1:8785
```

```
transport = sse

url = http://127.0.0.1:8785/
```

---

## Transport streamable HTTP

Сборка:

```bash
go build -o ./build/mcp-bitrix24-streamable ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-streamable
```

Запуск:

```bash
./build/mcp-bitrix24-streamable -listen 127.0.0.1:8786
```

```
transport = streamable

url = http://127.0.0.1:8786/
```

---

## Инструменты MCP

Схемы аргументов отдаёт сам сервер MCP (поля инструментов):

| Tool                                  | Назначение                                                                                                                                 |
|---------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| `b24_list_tasks`                      | `tasks.task.list` (`filter`, `select`, `order`, `params`, `start`)                                                                         |
| `b24_get_task`                        | `tasks.task.get` (`task_id`, `select`)                                                                                                     |
| `b24_get_task_comments`               | `task.commentitem.getlist` (`task_id`, `order`, `filter`)                                                                                  |
| `b24_get_task_timeline`               | Таймлайн задачи: события по задаче + комментарии (`task_id`, `include_comments`, `limit`)                                                  |
| `b24_analyze_task_blockers`           | Анализ блокеров по комментариям задачи (`task_id`, `limit`)                                                                                |
| `b24_analyze_task_execution_drift`    | Анализ дрифта исполнения задачи: факт/план и тишина коммуникаций (`task_id`)                                                               |
| `b24_analyze_responsible_performance` | Сводка по ответственному: объем, просрочки, high-risk, блокеры (`responsible_id`, `filter`, `order`, `start`, `limit`, `include_comments`) |
| `b24_analyze_project_health`          | Health-score портфеля задач и драйверы риска (`filter`, `order`, `start`, `limit`, `include_comments`)                                     |
| `b24_analyze_task`                    | Глубокий анализ одной задачи (`task_id`, `include_comments`)                                                                               |
| `b24_analyze_tasks_by_query`          | Аналитика по текстовому запросу (`query`, `task_id`, `filter`, `order`, `start`, `limit`, `include_comments`)                              |
| `b24_analyze_tasks_portfolio`         | Портфельная аналитика (`filter`, `order`, `start`, `limit`, `include_comments`, `group_by`)                                                |
| `b24_analyze_tasks_executive_summary` | Управленческая сводка за период (`filter`, `order`, `start`, `limit`, `period_days`, `include_comments`)                                   |
| `b24_analyze_tasks_sla`               | SLA-контроль (`filter`, `order`, `start`, `limit`, `soon_hours_threshold`, `include_comments`)                                             |
| `b24_analyze_tasks_workload`          | Баланс нагрузки по ответственным (`filter`, `order`, `start`, `limit`, `include_comments`, `overload_tasks`)                               |
| `b24_analyze_tasks_status_trends`     | Тренды по статусам (`filter`, `order`, `start`, `limit`, `period_days`)                                                                    |

## Mock REST (локальная отладка)

Отдельный бинарник - не MCP, а простой HTTP-мок Bitrix REST:

```bash
go build -o ./build/mcp-bitrix24-mock-rest ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-mock-rest
./build/mcp-bitrix24-mock-rest -listen 127.0.0.1:8899
```

База методов: `http://127.0.0.1:8899/rest/1/mock-token/<method>`.
