Инструкция сделана частично с помощью ии. Если будут вопросы по backend пишите.

# Подключение к backend

Base URL:
https://ai-cv-analyzer-production.up.railway.app

---

## Endpoints

| Метод | URL                       | Описание                                                 |
| ----- | ------------------------- | -------------------------------------------------------- |
| POST  | `/candidates`             | Создать кандидата (автоматически запускается AI скоринг) |
| GET   | `/candidates`             | Получить всех кандидатов                                 |
| GET   | `/candidates/:id`         | Получить одного кандидата                                |
| POST  | `/candidates/:id/score`   | Запустить AI оценку вручную                              |
| POST  | `/candidates/:id/approve` | Одобрить кандидата                                       |
| POST  | `/candidates/:id/reject`  | Отклонить кандидата                                      |
| GET   | `/leaderboard`            | Топ-10 кандидатов по score                               |

---

## Пример ответа кандидата

```json
{
  "id": 1,
  "name": "Арман",
  "age": 18,
  "status": "pending",
  "score": null
}
```

Статусы:
pending → scored → approved / rejected

---

## AI (Python)

Подними FastAPI на порту 8000

Endpoint:
POST /score

Backend отправляет:

```json
{
  "id": 1,
  "name": "...",
  "age": 18,
  "essay": "...",
  "experience": "...",
  "motivation": "..."
}
```

AI должен вернуть строго:

```json
{
  "score": 85,
  "explanation": "...",
  "ai_detected": false
}
```

После деплоя своего сервиса — скинь URL, я подключу его в backend.


## Фронтенд (React)

сценарий:

1. POST /candidates — создать кандидата (скоринг запускается автоматически)
2. GET /candidates — список кандидатов
3. GET /candidates/:id — детальная страница
4. POST /candidates/:id/approve или /reject — действия комиссии
5. GET /leaderboard — топ кандидатов

---

## Создание кандидата

POST /candidates

```json
{
  "name": "Арман",
  "age": 18,
  "essay": "текст эссе...",
  "experience": "опыт...",
  "motivation": "мотивация..."
}
```
