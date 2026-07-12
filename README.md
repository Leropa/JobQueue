# 🚀 Distributed Job Queue in Go

Асинхронная распределенная очередь задач с фоновым обработчиком (воркером), построенная на архитектуре **Clean Architecture** с использованием **Go**, **Redis** и **Docker**.

## 🏗️ Архитектура проекта

Проект спроектирован по канонам Чистой Архитектуры, что обеспечивает слабую связанность компонентов и легкую замену инфраструктуры (например, замену Redis на другую БД без изменения бизнес-логики).

*   `cmd/main.go` — Точка входа, сборка слоев («матрешка») и запуск HTTP-сервера.
*   `internal/domain/` — Бизнес-модели (структура Task) и абстракции (интерфейсы репозиториев).
*   `internal/repository/redis/` — Реализация работы с базой данных Redis (методы `Enqueue` и `Dequeue`).
*   `internal/usecase/` — Бизнес-логика, управление очередью и фоновый процесс воркера.
*   `internal/delivery/http/` — Транспортный слой (HTTP-хендлеры) для взаимодействия с внешним миром.

## 🛠️ Стек технологий

*   **Язык:** Go (Golang)
*   **База данных:** Redis (используется как FIFO-очередь с методами LPUSH/RPOP)
*   **Роутер:** [go-chi/chi](https://github.com/go-chi/chi) (v5)
*   **Драйвер БД:** [go-redis](https://github.com/redis/go-redis) (v9)
*   **Контейнеризация:** Docker + WSL2

## 🚀 Быстрый запуск

### 1. Запуск Redis
Убедитесь, что у вас запущен Docker Desktop, и поднимите контейнер с Redis одной командой:

```bash
docker run -d --name redis-jobqueue -p 6379:6379 redis:alpine
```

2. Запуск приложения
В корневой директории проекта выполните команду для старта сервера:

```bash
go run cmd/main.go
```

Вы увидите логи о том, что воркер успешно стартовал, а сервер слушает порт 8080.

📡 Эндпоинты API
Создание новой задачи
URL: http://localhost:8080/task

Метод: POST

Заголовки: Content-Type: application/json

Тело запроса (JSON):

```JSON
{
    "id": "task-777",
    "type": "send_email",
    "payload": "Welcome email to user@example.com"
}
```

Пример запроса через cURL:

```bash
curl -X POST http://localhost:8080/task \
     -H "Content-Type: application/json" \
     -d '{"id":"task-777","type":"send_email","payload":"Welcome email to user@example.com"}'
```
Успешный ответ (201 Created):

JSON
{
    "message": "Task created successfully"
}
🔄 Как это работает внутри
Клиент делает POST-запрос на /task.

TaskHandler парсит JSON и передает его в TaskUsecase.

TaskRepository сериализует структуру в строку и кладет ее в Redis с помощью LPUSH под ключом tasks_queue.

В этот же момент в фоне параллельно крутится горутина воркера (StartWorker), которая каждую секунду опрашивает Redis методом RPOP.

Как только в Redis появляется задача, воркер мгновенно забирает её, десериализует обратно в структуру Go и запускает «обработку» (выводит логи в консоль).
