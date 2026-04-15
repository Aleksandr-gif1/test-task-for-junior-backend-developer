# Task Service

REST API на Go для управления задачами с поддержкой **периодических (рекуррентных) задач**. Проект выполнен как тестовое задание — требования были описаны намеренно не детально, чтобы продемонстрировать способность самостоятельно проектировать решение.

## Что умеет

| Возможность | Описание |
|---|---|
| **CRUD задач-шаблонов** | Создание, чтение, обновление, удаление задач |
| **Периодические задачи** | Задача становится шаблоном — экземпляры создаются автоматически по расписанию |
| **4 типа расписания** | `daily`, `monthly_date`, `specific_dates`, `even_odd` |
| **Превентивная генерация** | Экземпляры физически записываются в БД — у каждого есть статус, который можно менять |
| **Фоновый scheduler** | Каждые 24 часа генерирует экземпляры на 7 дней вперёд |
| **Swagger-документация** | Встроенный UI + OpenAPI спецификация |

## Архитектура

Проект следует чистой слоистой структуре:

```
transport/http   →  usecase/task   →  repository/postgres   →  domain/task
   (маршруты)        (бизнес-           (хранение)              (модели)
                      логика)
```

- **Domain** — модели `Task` (шаблон) и `TaskInstance` (конкретный экземпляр на дату), статусы, типы периодичности
- **Use-case** — сервисная логика: CRUD, валидация, генератор дат по 4 алгоритмам
- **Repository** — работа с PostgreSQL через `pgx/v5` (connection pool)
- **Transport** — HTTP API на `gorilla/mux`, Swagger UI
- **Scheduler** — фоновый goroutine, запускающий генерацию по ticker'у

## Быстрый запуск

```bash
docker compose up --build
```

Сервис доступен на **http://localhost:8080**

> Если ранее запускался PostgreSQL — пересоздайте volume, чтобы применились миграции:
> ```bash
> docker compose down -v && docker compose up --build
> ```

### Ручной запуск (без Docker)

1. Создать БД `taskservice` в PostgreSQL
2. Применить миграции из `migrations/`
3. Указать DSN и запустить:

```bash
DATABASE_DSN="postgres://user:pass@localhost:5432/taskservice?sslmode=disable" go run ./cmd/api/
```

## API

Базовый префикс: `/api/v1`

### Шаблоны задач

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/tasks` | Создать задачу |
| `GET` | `/tasks` | Список всех задач |
| `GET` | `/tasks/{id}` | Получить по ID |
| `PUT` | `/tasks/{id}` | Обновить задачу |
| `DELETE` | `/tasks/{id}` | Удалить задачу (CASCADE → все экземпляры) |
| `POST` | `/tasks/generate?days=7` | Сгенерировать экземпляры вручную |

### Экземпляры

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/tasks/instances?from=2026-04-14&to=2026-04-21` | Экземпляры за период |
| `GET` | `/tasks/instances/{id}` | Конкретный экземпляр |
| `PUT` | `/tasks/instances/{id}/status` | Обновить статус (`new`, `in_progress`, `done`) |

### Swagger

- UI: http://localhost:8080/swagger/
- OpenAPI JSON: http://localhost:8080/swagger/openapi.json

## Периодические задачи

При создании/обновлении задачи можно указать `recurrence_type` и `recurrence_config`. Задача становится **шаблоном** (`is_recurring: true`), а scheduler автоматически создаёт экземпляры.

### Типы расписания

| Тип | `recurrence_config` | Логика |
|---|---|---|
| `daily` | `{"every": 2}` | Каждый N-й день от `created_at` шаблона |
| `monthly_date` | `{"day_of_month": 15}` | 15-го числа каждого месяца; если дня нет (31 в феврале) — берётся последний день |
| `specific_dates` | `{"dates": ["2026-04-20", "2026-05-10"]}` | Только явно указанные даты, попадающие в диапазон генерации |
| `even_odd` | `{"parity": "even"}` | Все чётные (`even`) или нечётные (`odd`) дни месяца |

### Примеры

```bash
# Ежедневная задача
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Обзвон пациентов","recurrence_type":"daily","recurrence_config":{"every":1}}'

# Сгенерировать экземпляры на 30 дней
curl -X POST "http://localhost:8080/api/v1/tasks/generate?days=30"

# Получить экземпляры за неделю
curl "http://localhost:8080/api/v1/tasks/instances?from=2026-04-14&to=2026-04-21"

# Отметить выполненным
curl -X PUT http://localhost:8080/api/v1/tasks/instances/42/status \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'
```

## Как тестировать

### 1. Запустить сервис

```bash
cd "/home/alex/Рабочий стол/test-task-for-junior-backend-developer"
docker compose up --build
```

### 2. Проверить здоровье

```bash
curl http://localhost:8080/api/v1/tasks
# [] — пусто, но 200 OK
```

### 3. Полный сценарий

```bash
# Создание периодического шаблона
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Ежедневный отчёт","recurrence_type":"daily","recurrence_config":{"every":1}}'

# Генерация экземпляров
curl -X POST http://localhost:8080/api/v1/tasks/generate

# Просмотр созданных экземпляров
curl "http://localhost:8080/api/v1/tasks/instances?from=2026-04-14&to=2026-04-21"

# Изменение статуса
curl -X PUT http://localhost:8080/api/v1/tasks/instances/1/status \
  -H "Content-Type: application/json" -d '{"status":"done"}'
```

### 4. Протестировать разные типы расписания

```bash
# monthly_date — каждый месяц 15-го числа
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Бухгалтерский отчёт","recurrence_type":"monthly_date","recurrence_config":{"day_of_month":15}}'

# even_odd — только чётные дни
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Инвентаризация","recurrence_type":"even_odd","recurrence_config":{"parity":"even"}}'

# specific_dates — конкретные даты
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Собрание","recurrence_type":"specific_dates","recurrence_config":{"dates":["2026-04-20","2026-05-10","2026-06-01"]}}'
```

## Переменные окружения

| Переменная | По умолчанию | Описание |
|---|---|---|
| `HTTP_ADDR` | `:8080` | Адрес HTTP-сервера |
| `DATABASE_DSN` | `postgres://postgres:postgres@localhost:5432/taskservice?sslmode=disable` | Строка подключения к PostgreSQL |
| `SCHEDULER_ENABLED` | `true` | Включить/выключить фоновый scheduler |

## Схема БД

```
tasks                          task_instances
┌────────────────────┐         ┌────────────────────────┐
│ id          BIGSER │         │ id              BIGSER │
│ title       TEXT   │  1───N  │ parent_task_id  BIGINT │ FK → tasks(id) CASCADE
│ description TEXT   │         │ scheduled_date  DATE   │
│ status      TEXT   │         │ title           TEXT   │
│ recurrence_type TEXT│        │ description     TEXT   │
│ recurrence_config JSONB│     │ status          TEXT   │
│ recurrence_active BOOL│      │ created_at      TIMEST │
│ created_at    TIMEST │       │ updated_at      TIMEST │
│ updated_at    TIMEST │       └────────────────────────┘
└────────────────────┘           UNIQUE (parent_task_id, scheduled_date)
```

## Принятые решения и допущения

### Модель «шаблон → экземпляр»

Периодическая задача — **шаблон** с настройками расписания. Конкретные задачи на каждый день — **экземпляры** в отдельной таблице `task_instances`. У каждого экземпляра свой независимый статус. Удаление шаблона каскадно удаляет все его экземпляры.

### Превентивная генерация в БД

Экземпляры создаются физически, а не генерируются on-the-fly при запросе. Это позволяет:
- Трекать статусы (выполнено/не выполнено)
- Видеть «план» задач заранее
- Фильтровать и искать по экземплярам

### Защита от дублей

`UNIQUE (parent_task_id, scheduled_date)` на уровне БД + проверка `Exists()` перед вставкой. Даже при двойном запуске scheduler'а дубли невозможны.

### Обновление шаблона не меняет экземпляры

При изменении расписания уже созданные экземпляры **не пересчитываются** — они уже в БД. Изменения применятся при следующей генерации для будущих дат.

### Всё в UTC

Временные пояса не поддерживаются — всё хранится и обрабатывается в UTC. Для продакшена можно добавить поле `timezone` в шаблон.

### Не реализовано (и почему)

| Что | Почему |
|---|---|
| `weekly` (еженедельно) | Не входило в 4 запрошенных типа, но архитектура расширяема — достаточно добавить новый `RecurrenceType` и алгоритм |
| Дата окончания периодичности | Можно добавить поле `recurrence_end`; scheduler будет пропускать шаблоны после этой даты |
| Уведомления/напоминания | Выходят за рамки CRUD + scheduling |
| Авторизация/роли | Не требуется по ТЗ; добавляется middleware в production |
| Тесты | Unit-тесты для `generator` и `validator` легко покрываются — структура позволяет добавить без рефакторинга |

## Стек

- **Go 1.23** — стандартная библиотека + `slog` для логирования
- **gorilla/mux** — HTTP-роутер
- **pgx/v5** — драйвер PostgreSQL с connection pool
- **Docker Compose** — оркестрация + init-миграции через `docker-entrypoint-initdb.d`
