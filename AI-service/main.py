"""
Gate & Crown AI Service
-----------------------
Provides three AI techniques:
  1. Vector Search  — semantic search over catalog using sentence-transformers + pgvector
  2. LLM Assistant  — structured order configuration help via Claude (Anthropic)
  3. RAG             — catalog/pricing Q&A with retrieved context fed into Claude
"""

import os
import json
from contextlib import asynccontextmanager
from typing import Optional

import anthropic
import numpy as np
import psycopg2
from dotenv import load_dotenv
from fastembed import TextEmbedding
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pgvector.psycopg2 import register_vector
from pydantic import BaseModel

load_dotenv()

DATABASE_URL = os.getenv("DATABASE_URL", "")
ANTHROPIC_API_KEY = os.getenv("ANTHROPIC_API_KEY", "")
PORT = int(os.getenv("AI_SERVICE_PORT", "8001"))

# ── globals ──────────────────────────────────────────────────────────────────
_model: Optional[TextEmbedding] = None
_db_conn = None
_anthropic_client: Optional[anthropic.Anthropic] = None


def get_model() -> TextEmbedding:
    global _model
    if _model is None:
        _model = TextEmbedding("BAAI/bge-small-en-v1.5")
    return _model


def get_db():
    global _db_conn
    if _db_conn is None or _db_conn.closed:
        _db_conn = psycopg2.connect(DATABASE_URL)
        register_vector(_db_conn)
    return _db_conn


def get_anthropic() -> anthropic.Anthropic:
    global _anthropic_client
    if _anthropic_client is None:
        _anthropic_client = anthropic.Anthropic(api_key=ANTHROPIC_API_KEY)
    return _anthropic_client


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Warm up embedding model at startup
    get_model()
    print("✅ Embedding model loaded (BAAI/bge-small-en-v1.5)")
    yield


app = FastAPI(title="Gate & Crown AI Service", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


# ── schemas ───────────────────────────────────────────────────────────────────

class SearchRequest(BaseModel):
    query: str
    category: Optional[str] = None
    limit: int = 10


class IngestRequest(BaseModel):
    batch_size: int = 50


class AssistRequest(BaseModel):
    user_message: str
    category: Optional[str] = None


class RagRequest(BaseModel):
    question: str


# ── helpers ───────────────────────────────────────────────────────────────────

def embed(text: str) -> list[float]:
    vec = list(get_model().embed([text]))[0]
    return vec.tolist()


def fetch_products_without_embeddings(cursor, batch_size: int):
    cursor.execute(
        "SELECT product_id, name, category, description FROM catalog_product "
        "WHERE embedding IS NULL AND status = 'active' LIMIT %s",
        (batch_size,),
    )
    return cursor.fetchall()


# ── routes ────────────────────────────────────────────────────────────────────

@app.get("/health")
def health():
    return {"status": "ok", "service": "gate-crown-ai"}


# ── AI Technique 1: Vector Search ─────────────────────────────────────────────

@app.post("/search")
def semantic_search(req: SearchRequest):
    """Semantic search over catalog using pgvector cosine similarity."""
    query_vec = embed(req.query)

    conn = get_db()
    cur = conn.cursor()

    sql = """
        SELECT product_id,
               1 - (embedding <=> %s::vector) AS similarity
        FROM catalog_product
        WHERE status = 'active'
          AND embedding IS NOT NULL
    """
    params: list = [query_vec]

    if req.category:
        sql += " AND category = %s"
        params.append(req.category)

    sql += " ORDER BY embedding <=> %s::vector LIMIT %s"
    params += [query_vec, req.limit]

    cur.execute(sql, params)
    rows = cur.fetchall()
    cur.close()

    return {"results": [{"product_id": r[0], "similarity": round(float(r[1]), 4)} for r in rows]}


# ── Data pipeline: ingest embeddings ─────────────────────────────────────────

@app.post("/pipeline/ingest")
def ingest_embeddings(req: IngestRequest):
    """Generate and store embeddings for catalog products that don't have one yet."""
    conn = get_db()
    cur = conn.cursor()

    rows = fetch_products_without_embeddings(cur, req.batch_size)
    if not rows:
        cur.close()
        return {"message": "All products already have embeddings.", "processed": 0}

    model = get_model()
    processed = 0
    for product_id, name, category, description in rows:
        text = f"{name}. {category}. {description or ''}"
        vec = list(model.embed([text]))[0].tolist()
        cur.execute(
            "UPDATE catalog_product SET embedding = %s::vector WHERE product_id = %s",
            (vec, product_id),
        )
        processed += 1

    conn.commit()
    cur.close()
    return {"message": f"Ingested {processed} products.", "processed": processed}


@app.get("/pipeline/status")
def pipeline_status():
    """Returns how many products have embeddings vs total."""
    conn = get_db()
    cur = conn.cursor()
    cur.execute("SELECT COUNT(*) FROM catalog_product WHERE status = 'active'")
    total = cur.fetchone()[0]
    cur.execute("SELECT COUNT(*) FROM catalog_product WHERE status = 'active' AND embedding IS NOT NULL")
    embedded = cur.fetchone()[0]
    cur.close()
    return {"total": total, "embedded": embedded, "pending": total - embedded}


# ── AI Technique 2: LLM Assistant (Claude) ────────────────────────────────────

@app.post("/assist")
def order_assistant(req: AssistRequest):
    """
    Uses Claude to help users configure gate/crown orders.
    Returns structured JSON with product recommendations and configuration.
    """
    client = get_anthropic()

    system_prompt = """You are an expert gate and crown product configurator for Gate & Crown B2B platform.
Help users configure orders by recommending suitable products and specifications.

Always respond with valid JSON in this exact structure:
{
  "recommendation": "brief recommendation text",
  "suggested_category": "frame|crown|mdf|glass|fittings|cutting",
  "key_features": ["feature1", "feature2"],
  "estimated_price_range": "$X - $Y",
  "next_steps": ["step1", "step2"]
}"""

    category_hint = f" Focus on {req.category} products." if req.category else ""

    message = client.messages.create(
        model="claude-haiku-4-5-20251001",
        max_tokens=512,
        system=system_prompt,
        messages=[{"role": "user", "content": req.user_message + category_hint}],
    )

    raw = message.content[0].text.strip()
    try:
        parsed = json.loads(raw)
    except json.JSONDecodeError:
        # Extract JSON from response if wrapped in markdown
        import re
        match = re.search(r"\{.*\}", raw, re.DOTALL)
        parsed = json.loads(match.group()) if match else {"recommendation": raw}

    return {"data": parsed, "model": "claude-haiku-4-5-20251001"}




if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=PORT, reload=True)
