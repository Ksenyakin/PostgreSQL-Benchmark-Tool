# Go Project: PostgreSQL Benchmark Tool

## Описание проекта

Это небольшой инструмент для тестирования производительности PostgreSQL, написанный на Go. Проект поддерживает многопоточные запросы к базе данных с настройками через JSON и предоставляет API для запуска бенчмарка через HTTP.

## Особенности

- Поддержка многопоточного выполнения SQL-запросов для имитации нагрузки на базу данных.
- Настройка параметров теста (подключение к базе данных, запрос, длительность теста, уровень параллелизма) через JSON.
- Поддержка API: возможность запуска бенчмарка с параметрами через HTTP-запрос.
- Логирование результатов и возможность их сохранения в файл.

## Установка

1. **Клонирование репозитория**:
   ```bash
   git clone https://github.com/ksenyakin/postgresql-benchmark-tool.git
   cd postgresql-benchmark-tool
   ```

2. **Список команд**:
- `make build`: сборка приложения Go.
- `make run-app`:  запуск только сервиса с приложением, игнорируя сервис PostgreSQL.
- `make docker-build`: Сборка приложения через Docker Compose..
- `make docker-run`: Запуск всех сервисов через Docker Compose.
- `make clean`:  Удаление собранных файлов и остановка всех контейнеров.
- `make deps`: Установка зависимостей.
- `make reset ` : Очистка томов и перезапуск контейнеров

## Структура проекта

```
.
.
├── cmd
│   └── app               // Точка входа в приложение (main.go)
├── internal
│   ├── application       // Логика приложения и сервисы
│   │   └── benchmark_service.go
│   ├── domain            // Доменные модели и интерфейсы
│   │   └── models.go
│   ├── infrastructure
│   │   ├── api           // API-интерфейс HTTP
│   │   │   └── handler.go
│   │   └── db            // Подключение к базе данных
│   │       └── postgres.go
│   └── repository        // Репозиторий для работы с данными бенчмарка
│       ├── benchmark_repository.go
│       ├── result_saver.go
│       └── workerpool.go
├── Dockerfile            // Dockerfile для контейнеризации
└── README.md             // Документация проекта

```

## Использование

### Конфигурация

Конфигурация задаётся в формате JSON и может быть передана через HTTP-запрос. Пример конфигурации:

```json
{
   "dsn": "postgres://postgres:12345678@postgres:5432/TestDB?sslmode=disable",
   "query": "SELECT * FROM users WHERE status  = 'active'",
   "duration_ms": 5000,
   "concurrency": 0
}
```

### Запуск через HTTP API

1. **Запуск сервера**:
   ```bash
   make docker-run
   ```
   Сервер запускается на `localhost:8080` и принимает запросы на эндпоинт `/benchmark`.


2. **Запуск бенчмарка через API**:
   Отправьте POST-запрос на `http://localhost:8080/benchmark` с JSON-конфигурацией в теле запроса.


3. **Пример использования через Postman**:
      - URL: `http://localhost:8080/benchmark`
      - Метод: `POST`
      - Тело (JSON):
      ```json
      {
         "dsn": "postgres://postgres:12345678@postgres:5432/TestDB?sslmode=disable",
         "query": "SELECT * FROM users WHERE status  = 'active'",
         "duration_ms": 5000,
         "concurrency": 0
      }
   ```
