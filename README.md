# Gen

Серверная и клиентская часть проекта Gen

---

### Клиентское приложение

| Платформа | Версия                                |
|-----------|---------------------------------------|
| Linux     | glibc 2.31+ (Ubuntu 20.04+ и аналоги) |
| Android   | 7.0+                                  |
| iOS       | 13.0+                                 |
| macOS     | Catalina 10.15+                       |
| Windows   | 10+                                   |

---

## Связанные репозитории

- **[gen-runner](https://github.com/magomedcoder/gen-runner)** - сервис запуска и взаимодействия с LLM

---

### Зависимости

- **Go** 1.26+
- **PostgreSQL** 16+
- **Клиент (Flutter/Dart):**
    - Flutter 3.24+
    - Dart SDK ^3.10.7
- **Protobuf** 30.2+
- 
---

### Сборка клиента

#### Linux и Android

Сборка **Linux** и **Android** через Docker:

```bash
docker build -f Dockerfile-client-app --target linux-build -t gen-app-linux .
docker run --rm -e TARGETS=linux,android -v ./out:/opt/gen/out gen-app-linux
```

#### Windows

Сборка **Windows** возможна только на **хосте Windows**:

```bash
docker build -f Dockerfile-client-app-windows --target windows-build -t gen-app-windows .
docker run --rm -v .\out:C:\gen\out gen-app-windows
```
