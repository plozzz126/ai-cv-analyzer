Сделал с помощью ии, если будут вопросы пишите

# Подключение к backend

Base URL: `http://localhost:8080`

Endpoints

| Метод | URL | Что делает |
| POST | `/candidates` | Создать заявку |
| GET | `/candidates` | Все кандидаты |
| GET | `/candidates/:id` | Один кандидат |
| POST | `/candidates/:id/score` | Запустить AI оценку |
| POST | `/candidates/:id/approve` | Одобрить |
| POST | `/candidates/:id/reject` | Отклонить |
| GET | `/leaderboard` | Топ-10 по score |


AI (Python)

Подними FastAPI на порту **8000** с одним endpoint:

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

---

## Фронтендер (React)

Типичный флоу:
```
1. POST /candidates → получаешь id
2. POST /candidates/:id/score → AI оценивает
3. GET /candidates → показываешь дашборд
4. POST /candidates/:id/approve или /reject → кнопки комиссии
```

Поля кандидата: `name, age, essay, experience, motivation` — все обязательные.

Статусы: `pending → scored → approved / rejected`
