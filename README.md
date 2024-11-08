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

2. **Установка зависимостей**:
   Убедитесь, что все зависимости загружены:
   ```bash
   go mod tidy
   ```

3. **Запуск**:
   ```bash
   go run .
   ```

## Структура проекта

```
.
├── main.go         // Основной файл, запускающий сервер и API
├── config.go       // Функции загрузки конфигурации
├── benchmark.go    // Основная логика бенчмарка
├── result.go       // Сохранение результатов бенчмарка
└── README.md       // Документация
```

## Использование

### Конфигурация

Конфигурация задаётся в формате JSON и может быть передана через HTTP-запрос. Пример конфигурации:

```json
{
  "dsn": "postgres://username:password@localhost:5432/dbname?sslmode=disable",
  "query": "SELECT * FROM users WHERE status = 'active'",
  "duration_ms": 10000,
  "concurrency": 5
}
```

### Запуск через HTTP API

1. **Запуск сервера**:
   ```bash
   go run .
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
       "dsn": "postgres://username:password@localhost:5432/dbname?sslmode=disable",
       "query": "SELECT * FROM users WHERE status = 'active'",
       "duration_ms": 10000,
       "concurrency": 5
     }
     ```
