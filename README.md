# LLM Runner

```bash
# Установка необходимых зависимостей и клонирование llama.cpp
make deps

# Генерация proto
make gen

# Сборка libllama.a (без CUDA)
make build-llama

# Сборка libllama.a с поддержкой NVIDIA (CUDA)
make build-llama-cublas

# Запуск (CPU, без CUDA)
make run-cpu

# Запуск (GPU, NVIDIA CUDA)
make run-gpu

# Сборка бинарника (CPU)
make build-cpu

# Сборка бинарника (CUDA)
make build-gpu
```

```bash
./build/llm-runner serve
```

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

## Скачивание весов (Hugging Face)

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
