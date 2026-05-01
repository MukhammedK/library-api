# Library API

REST API для управления библиотекой книг. Написан на Go с использованием PostgreSQL.

## Технологии

- Go 1.24
- PostgreSQL 16
- Docker / Docker Compose
- GitHub Actions (CI/CD)

## Запуск через Docker

```bash
docker-compose up --build
```

Сервер запустится на `http://localhost:8080`

## Endpoints

| Метод | URL | Описание |
|-------|-----|----------|
| GET | /book | Получить все книги |
| POST | /book | Добавить книгу |
| PUT | /book/{id} | Обновить книгу |
| DELETE | /book/{id} | Удалить книгу |

## Примеры запросов

**Добавить книгу:**
```bash
curl -X POST http://localhost:8080/book \
  -H "Content-Type: application/json" \
  -d '{"title": "The Go Programming Language", "year": 2015, "genre": "Programming"}'
```

**Получить все книги:**
```bash
curl http://localhost:8080/book
```

**Обновить книгу:**
```bash
curl -X PUT http://localhost:8080/book/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Go in Action", "year": 2016, "genre": "Programming"}'
```

**Удалить книгу:**
```bash
curl -X DELETE http://localhost:8080/book/1
```

## Локальный запуск

1. Создай файл `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=твой_пароль
DB_NAME=library
```

2. Запусти PostgreSQL и создай таблицу:
```bash
psql -U postgres -d library -f init.sql
```

3. Запусти сервер:
```bash
go run main.go
```