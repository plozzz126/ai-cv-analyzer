# Инструкция для фронтендера

## Base URL (единственный нужный тебе)
```
https://ai-cv-analyzer-production.up.railway.app
```

---

## Все endpoints

### 1. Создать кандидата
```
POST /candidates
```
Body (все поля обязательные!):
```json
{
  "name": "Арман",        // string
  "age": 18,              // number
  "essay": "текст...",    // string
  "experience": "опыт...",// string
  "motivation": "..."     // string
}
```
Ответ — объект кандидата с `id`. Скоринг запускается **автоматически** в фоне, ждать не нужно.

---

### 2. Получить всех кандидатов (для дашборда)
```
GET /candidates
```
Ответ:
```json
[
  {
    "id": 1,
    "name": "Арман",
    "age": 18,
    "essay": "...",
    "experience": "...",
    "motivation": "...",
    "score": 85,
    "explanation": "Кандидат показывает высокий потенциал...",
    "ai_detected": false,
    "status": "scored",
    "created_at": "2026-03-27T12:25:30Z"
  }
]
```

---

### 3. Получить одного кандидата
```
GET /candidates/:id
```

---

### 4. Одобрить кандидата (кнопка комиссии)
```
POST /candidates/:id/approve
```

---

### 5. Отклонить кандидата (кнопка комиссии)
```
POST /candidates/:id/reject
```

---

### 6. Лидерборд топ-10
```
GET /leaderboard
```
Ответ:
```json
[
  { "id": 1, "name": "Арман", "score": 85, "status": "scored" }
]
```

---

## Статусы кандидата
| Статус | Значение |
|---|---|
| `pending` | Заявка создана, AI ещё не оценил |
| `scored` | AI оценил, можно смотреть score |
| `approved` | Комиссия одобрила |
| `rejected` | Комиссия отклонила |

---

## Пример кода на React

```javascript
const API = 'https://ai-cv-analyzer-production.up.railway.app'

// Создать кандидата
const create = async (form) => {
  const res = await fetch(`${API}/candidates`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(form)
  })
  return res.json() // { id, name, status: 'pending', ... }
}

// Получить всех
const getAll = async () => {
  const res = await fetch(`${API}/candidates`)
  return res.json()
}

// Одобрить
const approve = async (id) => {
  await fetch(`${API}/candidates/${id}/approve`, { method: 'POST' })
}

// Отклонить
const reject = async (id) => {
  await fetch(`${API}/candidates/${id}/reject`, { method: 'POST' })
}

// Лидерборд
const leaderboard = async () => {
  const res = await fetch(`${API}/leaderboard`)
  return res.json()
}
```

---

## Важно
- После создания кандидата статус сначала `pending` — подожди 10-15 секунд и сделай GET снова, статус станет `scored` и появится score
- Можно сделать polling: каждые 3 секунды делать GET /candidates/:id пока status !== 'pending'
- CORS уже настроен, можно делать запросы прямо из браузера