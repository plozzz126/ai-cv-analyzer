сделано от части ии если будут вопросы по беку пишите. 

# Подключение к backend

Base URL: `https://ai-cv-analyzer-production.up.railway.app`

## Endpoints

| Метод | URL | Что делает |
|---|---|---|
| POST | `/candidates` | Создать заявку |
| GET | `/candidates` | Все кандидаты |
| GET | `/candidates/:id` | Один кандидат |
| POST | `/candidates/:id/score` | Запустить AI оценку |
| POST | `/candidates/:id/approve` | Одобрить |
| POST | `/candidates/:id/reject` | Отклонить |
| GET | `/leaderboard` | Топ-10 по score |

---

## AI (Python)

Подними FastAPI на порту **8000**:

```
POST /score
```

Получишь:
```json
{ "id": 1, "name": "...", "age": 18, "essay": "...", "experience": "...", "motivation": "..." }
```

Верни строго:
```json
{ "score": 85, "explanation": "...", "ai_detected": false }
```

После деплоя своего сервиса — скинь мне URL, я добавлю в backend переменную AI_SERVICE_URL.

---

## Фронтендер (React)

Типичный флоу:
```
1. POST /candidates        → скоринг запускается автоматически
2. GET  /candidates        → дашборд (статус pending → scored)
3. GET  /candidates/:id    → детальная страница
4. POST /candidates/:id/approve или /reject → кнопки комиссии
5. GET  /leaderboard       → топ-10
```

Поля для POST /candidates:
```json
{
  "name": "Арман",
  "age": 18,
  "essay": "текст эссе...",
  "experience": "опыт...",
  "motivation": "мотивация..."
}
```

Статусы: `pending → scored → approved / rejected`