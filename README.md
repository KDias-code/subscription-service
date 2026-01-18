# 🟢 Subscriptions Service

Сервис для работы с подписками, реализованный на **Go** с **Fiber**, **PostgreSQL** и Docker.

> 🚀 Запуск через Docker Compose:
```bash
docker-compose up -d

swag init --generalInfo internal/handlers/subscruptionsHandlers.go --output docs

📁 Структура проекта

internal/handlers — HTTP-хендлеры (CRUD + подсчет суммы подписок)

internal/service — бизнес-логика

internal/model — модели данных

internal/repository — работа с базой данных

migrations — SQL миграции для PostgreSQL

Dockerfile — сборка приложения

docker-compose.yml — запуск приложения и базы данных

⚡ API Endpoints

Все эндпоинты находятся под префиксом /v1/subscriptions.

1️⃣ Получить подписку по ID
GET /v1/subscriptions/{id}


Параметры:

Параметр	Тип	Обязательный	Описание
id	        string	да	        UUID подписки

Ответы:

200 OK — объект подписки

400 Bad Request — некорректный UUID

404 Not Found — подписка не найдена

500 Internal Server Error — ошибка сервера

Пример:

curl -X GET http://localhost:8080/v1/subscriptions/uuid

2️⃣ Создать подписку
POST /v1/subscriptions


Тело запроса (JSON):

{
  "user_id": "uuid пользователя",
  "service_name": "название сервиса",
  "amount": 1000,
  "start_date": "2026-01-18",
  "end_date": "2026-02-18"
}


Ответ:

{
  "success": true
}


Пример curl:

curl -X POST http://localhost:8080/v1/subscriptions \
-H "Content-Type: application/json" \
-d '{"user_id":"uuid","service_name":"Test Service","amount":1000,"start_date":"2026-01-18","end_date":"2026-02-18"}'

3️⃣ Обновить подписку
PUT /v1/subscriptions


Тело запроса: та же структура, что для создания, но с существующим id.

Ответ:

{
  "success": true
}


Пример curl:

curl -X PUT http://localhost:8080/v1/subscriptions \
-H "Content-Type: application/json" \
-d '{"id":"uuid","user_id":"uuid","service_name":"Updated Service","amount":1200,"start_date":"2026-01-18","end_date":"2026-02-18"}'

4️⃣ Удалить подписку
DELETE /v1/subscriptions/{id}


Параметры:

Параметр	Тип	Обязательный	Описание
id	string	да	UUID подписки

Ответ:

{
  "success": true
}


Пример curl:

curl -X DELETE http://localhost:8080/v1/subscriptions/uuid

5️⃣ Получить сумму подписок по фильтрам
GET /v1/subscriptions/sum


Параметры запроса:

Параметр	Тип	Обязательный	Описание
startDate	string	да	Начало периода (YYYY-MM-DD)
endDate	string	да	Конец периода (YYYY-MM-DD)
userID	string	нет	UUID пользователя
serviceName	string	нет	Название сервиса

Пример запроса:

curl "http://localhost:8080/v1/subscriptions/sum?startDate=2026-01-01&endDate=2026-01-31&userID=uuid"


Ответ:

{
  "amount": 2500
}

🛠 Локальная разработка

Убедитесь, что Docker и Docker Compose установлены.

Запустите сервис и базу данных:

docker-compose up -d


Проверка подключения к БД:

docker exec -it subscriptions-postgres psql -U test -d subscriptions_db

📖 Логи и отладка

Логи сервиса:

docker logs -f subscriptions-service


Логи PostgreSQL:

docker logs -f subscriptions-postgres

🛠 Технологии

Go 1.25+

Fiber

PostgreSQL 16

Docker & Docker Compose

UUID, HCLog, CErrors