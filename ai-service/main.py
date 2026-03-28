from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy import create_engine, Column, Integer, String, Text, DateTime
from sqlalchemy.orm import declarative_base, sessionmaker, Session
from pydantic import BaseModel, field_validator
from groq import Groq
from dotenv import load_dotenv
from datetime import datetime, timezone
import os
import json
import logging
import re

# ================= LOGGING =================
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger(__name__)

# ================= ENV =================
load_dotenv()

GROQ_API_KEY = os.getenv("GROQ_API_KEY")
DB_URL = os.getenv("DATABASE_URL")
MODEL = os.getenv("GROQ_MODEL", "llama-3.3-70b-versatile")

if not GROQ_API_KEY:
    raise RuntimeError("GROQ_API_KEY is not set in environment variables")

# ================= DB =================
engine = create_engine(
    DB_URL,
    client_encoding="utf8",
    pool_pre_ping=True,
    pool_size=10,
    max_overflow=20,
)
SessionLocal = sessionmaker(bind=engine, autocommit=False, autoflush=False)
Base = declarative_base()


class Message(Base):
    __tablename__ = "messages"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(String(128), index=True, nullable=False)
    role = Column(String(32), nullable=False)
    content = Column(Text, nullable=False)
    created_at = Column(DateTime(timezone=True), default=lambda: datetime.now(timezone.utc))


Base.metadata.create_all(bind=engine)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# ================= APP =================
@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("Application startup")
    yield
    logger.info("Application shutdown")


app = FastAPI(
    title="AI Candidate Scoring API",
    version="2.0.0",
    lifespan=lifespan,
)

groq_client = Groq(api_key=GROQ_API_KEY)

# ================= PROMPT =================
SYSTEM_PROMPT = """
Ты — эксперт по отбору кандидатов.

Тебе приходит текст (или набор ответов) кандидата.

Оцени по четырём критериям (каждый от 0 до 25):
1. Мотивация
2. Лидерство
3. Опыт
4. Потенциал

Суммарный score = сумма четырёх оценок (0–100).

Также оцени вероятность использования ИИ при написании ответов (0–100).

Верни строго валидный JSON без каких-либо пояснений:

{
  "score": <int 0–100>,
  "recommendation": "<pass|maybe|reject>",
  "ai_probability": <int 0–100>,
  "strengths": ["<строка>", ...],
  "weaknesses": ["<строка>", ...],
  "explanation": "<строка>"
}
""".strip()

# ================= SCHEMAS =================
class CandidateRequest(BaseModel):
    user_id: str
    text: str

    @field_validator("user_id", "text")
    @classmethod
    def not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("Field must not be empty")
        return v.strip()


class AnalyzeResponse(BaseModel):
    user_id: str
    score: int
    recommendation: str
    ai_probability: int
    strengths: list[str]
    weaknesses: list[str]
    explanation: str


class StartRequest(BaseModel):
    user_id: str

    @field_validator("user_id")
    @classmethod
    def not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("user_id must not be empty")
        return v.strip()


class InterviewRequest(BaseModel):
    user_id: str
    answer: str

    @field_validator("user_id", "answer")
    @classmethod
    def not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("Field must not be empty")
        return v.strip()


class InterviewQuestionResponse(BaseModel):
    step: int
    total: int
    question: str


# ========== SCHEMA ДЛЯ GO BACKEND ==========
class ScoreRequest(BaseModel):
    id: int
    name: str
    age: int
    essay: str
    experience: str
    motivation: str


class ScoreResponse(BaseModel):
    score: int
    explanation: str
    ai_detected: bool


# ================= IN-MEMORY INTERVIEW STATE =================
interviews: dict[str, dict] = {}

QUESTIONS: list[str] = [
    "Расскажи о себе и своих целях.",
    "Какой у тебя был опыт лидерства?",
    "Расскажи про сложность, которую ты преодолел.",
    "Почему ты хочешь учиться в inVision U?",
]

# ================= UTILS =================
_JSON_RE = re.compile(r"\{.*\}", re.DOTALL)


def extract_json(text: str) -> dict:
    text = text.strip()
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        pass
    match = _JSON_RE.search(text)
    if match:
        try:
            return json.loads(match.group())
        except json.JSONDecodeError:
            pass
    fenced = re.sub(r"```(?:json)?", "", text).strip().strip("`").strip()
    return json.loads(fenced)


def call_groq(messages: list[dict]) -> tuple[dict, str]:
    completion = groq_client.chat.completions.create(
        model=MODEL,
        messages=messages,
        temperature=0.3,
        max_tokens=700,
    )
    reply = completion.choices[0].message.content
    logger.debug("Groq raw reply: %s", reply)
    return extract_json(reply), reply


def save_messages(db: Session, user_id: str, role: str, content: str) -> None:
    db.add(Message(user_id=user_id, role=role, content=content))
    db.commit()


# ================= ROUTES =================

@app.get("/health", tags=["System"])
def health():
    return {"status": "ok"}


# ========== ГЛАВНЫЙ ENDPOINT ДЛЯ GO BACKEND ==========
@app.post("/score", response_model=ScoreResponse, tags=["Scoring"])
def score(req: ScoreRequest, db: Session = Depends(get_db)):
    """
    Этот endpoint вызывает Go backend.
    Принимает данные кандидата, возвращает score, explanation, ai_detected.
    """
    text = f"""
Имя: {req.name}
Возраст: {req.age}

Эссе:
{req.essay}

Опыт:
{req.experience}

Мотивация:
{req.motivation}
""".strip()

    try:
        result, raw_reply = call_groq([
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": text},
        ])
        save_messages(db, str(req.id), "scoring", raw_reply)
    except json.JSONDecodeError as e:
        logger.error("JSON parse error: %s", e)
        raise HTTPException(status_code=502, detail=f"Model returned invalid JSON: {e}")
    except Exception as e:
        logger.exception("Groq API error")
        raise HTTPException(status_code=502, detail=str(e))

    ai_probability = result.get("ai_probability", 0)

    return ScoreResponse(
        score=result.get("score", 50),
        explanation=result.get("explanation", ""),
        ai_detected=ai_probability > 60,
    )


@app.post("/analyze", response_model=AnalyzeResponse, tags=["Analysis"])
def analyze(req: CandidateRequest, db: Session = Depends(get_db)):
    """Одиночный анализ текста кандидата."""
    try:
        result, raw_reply = call_groq([
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": req.text},
        ])
        save_messages(db, req.user_id, "analysis", raw_reply)
    except json.JSONDecodeError as e:
        logger.error("JSON parse error: %s", e)
        raise HTTPException(status_code=502, detail=f"Model returned invalid JSON: {e}")
    except Exception as e:
        logger.exception("Groq API error")
        raise HTTPException(status_code=502, detail=str(e))

    return AnalyzeResponse(user_id=req.user_id, **result)


@app.post("/interview/start", response_model=InterviewQuestionResponse, tags=["Interview"])
def start_interview(req: StartRequest):
    """Начать новое интервью."""
    interviews[req.user_id] = {"step": 0, "answers": []}
    logger.info("Interview started for user_id=%s", req.user_id)
    return InterviewQuestionResponse(step=1, total=len(QUESTIONS), question=QUESTIONS[0])


@app.post("/interview/answer", tags=["Interview"])
def answer_interview(req: InterviewRequest, db: Session = Depends(get_db)):
    """Ответить на текущий вопрос интервью."""
    interview = interviews.get(req.user_id)
    if not interview:
        raise HTTPException(status_code=404, detail="Interview not found. Call /interview/start first.")

    interview["answers"].append(req.answer)
    interview["step"] += 1

    if interview["step"] >= len(QUESTIONS):
        return _finish_interview(req.user_id, db)

    return InterviewQuestionResponse(
        step=interview["step"] + 1,
        total=len(QUESTIONS),
        question=QUESTIONS[interview["step"]],
    )


def _finish_interview(user_id: str, db: Session) -> dict:
    interview = interviews.pop(user_id)
    full_text = "\n\n".join(
        f"Вопрос {i+1}: {QUESTIONS[i]}\nОтвет: {ans}"
        for i, ans in enumerate(interview["answers"])
    )

    try:
        result, raw_reply = call_groq([
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": full_text},
        ])
        save_messages(db, user_id, "interview_answers", full_text)
        save_messages(db, user_id, "interview_analysis", raw_reply)
    except json.JSONDecodeError as e:
        logger.error("JSON parse error: %s", e)
        raise HTTPException(status_code=502, detail=f"Model returned invalid JSON: {e}")
    except Exception as e:
        logger.exception("Error finishing interview for user_id=%s", user_id)
        raise HTTPException(status_code=502, detail=str(e))

    logger.info("Interview finished for user_id=%s, score=%s", user_id, result.get("score"))

    return {
        "status": "finished",
        "user_id": user_id,
        **result,
    }