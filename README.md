# LLM Runner

## Быстрый старт (сборка и запуск)

```bash
# Установка зависимостей
make deps

# Генерация proto
make gen

# Сборка библиотек (без CUDA)
make build-libs-cpu

# Сборка библиотек с поддержкой NVIDIA (CUDA)
make build-libs-gpu

# Запуск (CPU, без CUDA)
make run-cpu

# Запуск (GPU, NVIDIA CUDA)
make run-gpu

# Сборка бинарника (CPU)
make build-cpu

# Сборка бинарника (CUDA)
make build-gpu
```

---

## Скачивание модели (Hugging Face)

```bash
./build/llm-runner download --repo <org/model> --list
./build/llm-runner download --repo <org/model> --file ....gguf
```

## Клиент к запущенному раннеру

```bash
./build/llm-runner remote ping
./build/llm-runner remote models
./build/llm-runner remote run --model mymodel --prompt "привет"
```

---

```bash
# Собрать yaml из Modelfile
./build/llm-runner create myalias -f ./Modelfile [--force]

# Показать yaml манифеста или экспорт в Modelfile
./build/llm-runner show myalias
./build/llm-runner show myalias --modelfile # или -m

# Только путь к .yaml (для скриптов)
./build/llm-runner show myalias --path-only

# Список локальных .gguf в каталоге model_path
./build/llm-runner models
```
