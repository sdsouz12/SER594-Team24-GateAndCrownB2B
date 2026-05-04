"""
AI-service unit tests — all external I/O (DB, Anthropic, embedding model) is mocked.
Run with: pytest tests/ -v  (from AI-service/)
"""

import numpy as np
import pytest
from unittest.mock import MagicMock, patch

# ── Pre-set the embedding model so the lifespan does not download it ──────────
import main as app_module

_mock_model = MagicMock()
_mock_model.embed.side_effect = lambda texts: iter([np.zeros(384)])
app_module._model = _mock_model

from fastapi.testclient import TestClient  # noqa: E402 (import after patching)

client = TestClient(app_module.app)


# ── Helpers ───────────────────────────────────────────────────────────────────

def _mock_db(fetchall=None, fetchone_seq=None):
    conn = MagicMock()
    cur = MagicMock()
    conn.cursor.return_value = cur
    if fetchall is not None:
        cur.fetchall.return_value = fetchall
    if fetchone_seq is not None:
        cur.fetchone.side_effect = fetchone_seq
    return conn


def _mock_ai(text: str):
    ai = MagicMock()
    msg = MagicMock()
    msg.content = [MagicMock(text=text)]
    ai.messages.create.return_value = msg
    return ai


# ── Health ────────────────────────────────────────────────────────────────────

def test_health_returns_ok():
    resp = client.get("/health")
    assert resp.status_code == 200
    assert resp.json()["status"] == "ok"
    assert resp.json()["service"] == "gate-crown-ai"


# ── Vector Search ─────────────────────────────────────────────────────────────

def test_search_missing_query_returns_422():
    resp = client.post("/search", json={})
    assert resp.status_code == 422


def test_search_returns_empty_list_when_no_results():
    db = _mock_db(fetchall=[])
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.post("/search", json={"query": "heavy steel gate"})
    assert resp.status_code == 200
    assert resp.json()["results"] == []


def test_search_returns_results_with_similarity():
    rows = [(1, 0.93), (2, 0.87), (3, 0.76)]
    db = _mock_db(fetchall=rows)
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.post("/search", json={"query": "sliding gate", "limit": 3})
    assert resp.status_code == 200
    results = resp.json()["results"]
    assert len(results) == 3
    assert results[0]["product_id"] == 1
    assert results[0]["similarity"] == 0.93


def test_search_with_category_filter():
    rows = [(5, 0.88)]
    db = _mock_db(fetchall=rows)
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.post("/search", json={"query": "ornate crown", "category": "crown"})
    assert resp.status_code == 200
    assert len(resp.json()["results"]) == 1


# ── Pipeline ──────────────────────────────────────────────────────────────────

def test_pipeline_status_returns_correct_counts():
    db = _mock_db(fetchone_seq=[(20,), (15,)])
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.get("/pipeline/status")
    assert resp.status_code == 200
    body = resp.json()
    assert body["total"] == 20
    assert body["embedded"] == 15
    assert body["pending"] == 5


def test_pipeline_ingest_no_products():
    db = _mock_db(fetchall=[])
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.post("/pipeline/ingest", json={})
    assert resp.status_code == 200
    assert resp.json()["processed"] == 0


def test_pipeline_ingest_processes_batch():
    rows = [
        (1, "Classic Frame Gate", "frame", "Traditional frame gate"),
        (2, "Crown Panel Standard", "crown", "Premium crown panel"),
    ]
    db = _mock_db(fetchall=rows)
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.post("/pipeline/ingest", json={"batch_size": 10})
    assert resp.status_code == 200
    assert resp.json()["processed"] == 2


# ── LLM Assistant ─────────────────────────────────────────────────────────────

def test_assist_missing_message_returns_422():
    resp = client.post("/assist", json={})
    assert resp.status_code == 422


def test_assist_returns_structured_json():
    payload = (
        '{"recommendation":"Use steel frame",'
        '"suggested_category":"frame",'
        '"key_features":["durable","weatherproof"],'
        '"estimated_price_range":"$800 - $1500",'
        '"next_steps":["measure opening","choose finish"]}'
    )
    ai = _mock_ai(payload)
    with patch.object(app_module, "get_anthropic", return_value=ai):
        resp = client.post("/assist", json={"user_message": "I need a large driveway gate"})
    assert resp.status_code == 200
    data = resp.json()["data"]
    assert "recommendation" in data
    assert data["suggested_category"] == "frame"


def test_assist_with_category_hint():
    payload = (
        '{"recommendation":"Crown panel suits your style",'
        '"suggested_category":"crown",'
        '"key_features":["ornate"],'
        '"estimated_price_range":"$300 - $600",'
        '"next_steps":["select size"]}'
    )
    ai = _mock_ai(payload)
    with patch.object(app_module, "get_anthropic", return_value=ai):
        resp = client.post("/assist", json={"user_message": "decorative top piece", "category": "crown"})
    assert resp.status_code == 200
    assert resp.json()["data"]["suggested_category"] == "crown"


def test_assist_handles_plain_text_response_gracefully():
    ai = _mock_ai("I recommend a steel gate with crown panel for your entrance.")
    with patch.object(app_module, "get_anthropic", return_value=ai):
        resp = client.post("/assist", json={"user_message": "help me choose"})
    assert resp.status_code == 200
    assert "recommendation" in resp.json()["data"]


# ── RAG ───────────────────────────────────────────────────────────────────────

def test_rag_missing_question_returns_422():
    resp = client.post("/rag", json={})
    assert resp.status_code == 422


def test_rag_returns_503_when_no_catalog():
    db = _mock_db(fetchall=[])
    with patch.object(app_module, "get_db", return_value=db):
        resp = client.post("/rag", json={"question": "What gates do you sell?"})
    assert resp.status_code == 503


def test_rag_returns_answer_and_sources():
    rows = [
        ("Classic Frame Gate", "frame", "Traditional gate", "$800-$1500", "Steel", "3m x 2m"),
        ("Crown Panel Standard", "crown", "Decorative crown", "$300-$600", "Iron", "1m x 0.5m"),
    ]
    db = _mock_db(fetchall=rows)
    ai = _mock_ai("The Classic Frame Gate is ideal for residential driveways.")
    with patch.object(app_module, "get_db", return_value=db), \
         patch.object(app_module, "get_anthropic", return_value=ai):
        resp = client.post("/rag", json={"question": "Which gate for a driveway?"})
    assert resp.status_code == 200
    body = resp.json()
    assert "answer" in body
    assert "sources" in body
    assert len(body["sources"]) == 2
    assert body["sources"][0] == "Classic Frame Gate"
