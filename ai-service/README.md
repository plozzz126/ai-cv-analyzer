# AI Candidate Scoring System

## 📌 Описание

AI-система для автоматизированного отбора кандидатов в inVision U (Decentrathon 5.0).

Система анализирует кандидатов через:

* текст (эссе)
* интерактивное интервью

И выдает:

* score (0–100)
* рекомендацию (pass / maybe / reject)
* вероятность использования AI
* объяснение (Explainable AI)

---

## 🚀 Функционал

### 🤖 Анализ кандидата

* Оценка мотивации, опыта, лидерства и потенциала
* AI detection

### 🎤 Интервью-бот

* Пошаговые вопросы
* Сбор ответов
* Финальный анализ

### 💾 Хранение

* PostgreSQL база данных
* История интервью и анализов

---

## 🛠️ Технологии

* FastAPI
* PostgreSQL
* SQLAlchemy
* Groq LLM (llama-3.3-70b)
* Pydantic

---

## ⚙️ Установка

### 1. Клонировать проект

```bash
git clone <repo_url>
cd project
```

### 2. Установить зависимости

```bash
pip install -r requirements.txt
```

### 3. Настроить .env

```env
GROQ_API_KEY=your_key
DATABASE_URL=postgresql://postgres:1234@localhost:5432/chatdb
```

### 4. Запуск

```bash
uvicorn main:app --reload
```

---

## 📡 API

### 🔹 Проверка

GET /health

### 🔹 Анализ

POST /analyze

```json
{
  "user_id": "123",
  "text": "..."
}
```

### 🔹 Интервью

#### Старт

POST /interview/start

```json
{
  "user_id": "123"
}
```

#### Ответ

POST /interview/answer

```json
{
  "user_id": "123",
  "answer": "..."
}
```

---

## 🧠 Архитектура

User → FastAPI → Groq LLM → Scoring → PostgreSQL

---

## 🎯 Цель

Сократить ручной отбор кандидатов и выявлять скрытый потенциал.

---

## ⚠️ Ограничения

* Интервью хранится в памяти (in-memory)
* При перезапуске данные интервью теряются
* Требуется PostgreSQL

---

## 🔮 Будущее развитие

* Redis для хранения интервью
* Веб-интерфейс для комиссии
* Улучшение AI detection
* Ranking кандидатов
