# Gate & Crown B2B Platform

**SER 594 — AI for Software Engineers | Team 24 | Final Milestone**

Full-stack B2B ordering platform with semantic search, LLM order assistant, and RAG — built with Go, Vue 3, and Python (FastAPI).

**CI Status:** [![CI](https://github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/actions/workflows/ci.yml/badge.svg)](https://github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/actions/workflows/ci.yml)
**CI Dashboard:** https://github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/actions

---

## What's Implemented (Feature Complete)

| Feature | Status |
|---|---|
| JWT authentication — login, registration, role-based sessions | ✅ Done |
| AI Technique 1 — Vector Search (BAAI/bge-small-en-v1.5 + pgvector) | ✅ Done |
| AI Technique 2 — LLM Order Assistant (Claude Haiku, structured JSON) | ✅ Done |
| AI Technique 3 — RAG (catalog context retrieval + Claude grounded Q&A) | ✅ Done |
| Order management — place orders from catalog, track status | ✅ Done |
| Test suite — 30 automated unit tests (Go + Python) | ✅ Done |
| Docker deployment — all services run with one command | ✅ Done |
| GitHub Actions CI — tests run on every push automatically | ✅ Done |

---

## CI — Already Running

GitHub Actions is configured at `.github/workflows/ci.yml` and runs automatically on every push to `main`. **No manual configuration needed.**

View live CI results:
**https://github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/actions**

---

## Option A — Docker (Recommended, Easiest)

No manual setup required — Docker handles the database, migrations, backend, AI service, and frontend automatically.

```bash
docker compose up --build
```

Then open **http://localhost:3000** and log in.

> First run downloads the embedding model (~66 MB inside the container) — allow ~2–3 minutes.

To stop:
```bash
docker compose down
```

To wipe data and start completely fresh:
```bash
docker compose down -v
docker compose up --build
```

---

## Option B — Manual Local Setup

### Step 1 — Install pgvector

```bash
brew install pgvector
psql postgres -c "CREATE EXTENSION IF NOT EXISTS vector;"
```

### Step 2 — Configure Backend Environment

```bash
cp backend-service/.env.example backend-service/.env
```

Edit `backend-service/.env` and set your PostgreSQL credentials:

```
DATABASE_URL=postgresql://localhost:5432/gate_crown
PORT=8080
JWT_SECRET=change-me-to-a-long-random-secret
JWT_EXPIRATION_HOURS=24
FRONTEND_URL=http://localhost:5173
AI_SERVICE_URL=http://localhost:8001
```

### Step 3 — Run Database Migrations

```bash
go run ./backend-service/cmd/migrate up
```

### Step 4 — Start the Backend

```bash
go run ./backend-service/cmd/server
```

Backend runs at: http://localhost:8080

### Step 5 — Set Up the AI Service

```bash
conda create -n gatecrown python=3.11 -y
conda activate gatecrown
```

The `.env.example` already contains the Anthropic API key — no changes needed:

```bash
cp AI-service/.env.example AI-service/.env
```

Install and start:

```bash
cd AI-service
/opt/anaconda3/envs/gatecrown/bin/pip install -r requirements.txt
/opt/anaconda3/envs/gatecrown/bin/python main.py
```

AI service runs at: http://localhost:8001

### Step 6 — Ingest Catalog Embeddings (once)

```bash
curl -X POST http://localhost:8001/pipeline/ingest \
  -H "Content-Type: application/json" -d '{}'
```

Expected: `{"message": "Ingested 20 products.", "processed": 20}`

### Step 7 — Start the Frontend

```bash
cd frontend-client
npm install
npm run dev
```

Frontend runs at: http://localhost:5173

---

## Using the Application

1. Open the app (http://localhost:5173 or http://localhost:3000 for Docker)
2. **Register** a new account or use a test account below
3. After login you land on the **Final Milestone Dashboard** showing all features
4. Go to **Catalog** to browse 20 products and use AI-powered search
5. Click **Place Order** on any product to submit an order
6. Go to **My Orders** to see your order history and status
7. Use the **Order Assistant** (Claude) on the catalog page for AI recommendations
8. Use the **RAG** section on the dashboard to ask questions about catalog pricing and materials

---

## Test Accounts

| Username | Password | Role |
|---|---|---|
| admin | admin123 | SUPERADMIN |
| john.doe | admin123 | ADMIN (Acme Corp) |
| jane.smith | admin123 | ADMIN (Acme Corp) |
| sarah.jones | admin123 | ADMIN (Beta Solutions) |

---

## Running the Test Suite

All 30 tests run **without a database or network connection** — all external I/O is mocked.

### Backend — 15 Go unit tests

```bash
cd backend-service
go test ./... -v
```

Covers: JWT validation (valid/expired/wrong secret/malformed), login (valid credentials/wrong password/deactivated/no role), register (success/username taken), profile update (no fields/short password/missing current/wrong current password).

### AI Service — 15 Python unit tests

```bash
cd AI-service
pip install -r requirements-test.txt
pytest tests/ -v
```

Covers: health endpoint, vector search (empty/results/category filter), pipeline status/ingest, LLM assistant (structured JSON/category hint/plain text fallback), RAG (503 when no catalog/success with sources).

### CI (automatic)

Both suites run automatically on every push to `main` via GitHub Actions.
Results: **https://github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/actions**

---

## AI Techniques Implemented

| Technique | How it works | Endpoint |
|---|---|---|
| Vector Search | User query embedded with BAAI/bge-small-en-v1.5 (ONNX/fastembed), compared against stored product embeddings in PostgreSQL using pgvector cosine similarity | `POST /api/catalog/search` |
| LLM Assistant | User describes a requirement; Claude Haiku returns structured JSON with product category, price range, key features, and next steps | `POST /api/ai/assist` |
| RAG | Top-5 most relevant products retrieved from pgvector and injected as context; Claude generates a grounded natural-language answer | `POST /api/ai/rag` |

---

## Project Structure

```
gate-crown/
├── backend-service/          Go + Gin REST API
│   ├── cmd/server/           Server entry point
│   ├── cmd/migrate/          Database migration tool
│   ├── internal/auth/        Auth (login, register, JWT) + 15 unit tests
│   ├── internal/catalog/     Catalog (list + semantic search)
│   ├── internal/orders/      Order management (create, list)
│   ├── internal/aiproxy/     Proxy to Python AI service
│   ├── migrations/           SQL migration files (001–007)
│   ├── Dockerfile
│   └── entrypoint.sh
├── frontend-client/          Vue 3 + Vite + Tailwind CSS
│   ├── src/views/
│   │   ├── WelcomeView.vue   Final milestone dashboard
│   │   ├── CatalogView.vue   Product catalog + AI search + order modal
│   │   ├── OrdersView.vue    Order history and status
│   │   ├── LoginView.vue
│   │   └── RegisterView.vue
│   ├── Dockerfile
│   └── nginx.conf
├── AI-service/               Python FastAPI AI microservice
│   ├── main.py               Vector search, LLM assistant, RAG
│   ├── tests/test_main.py    15 Python unit tests
│   ├── requirements.txt
│   ├── requirements-test.txt
│   └── Dockerfile
├── docker-compose.yml        One-command deployment
├── .github/workflows/ci.yml  GitHub Actions CI
└── go.work
```

---

## Quick Start (manual)

```bash
# 1. Migrations
go run ./backend-service/cmd/migrate up

# 2. Backend (terminal 1)
go run ./backend-service/cmd/server

# 3. AI service (terminal 2)
/opt/anaconda3/envs/gatecrown/bin/python AI-service/main.py

# 4. Ingest embeddings (once)
curl -X POST http://localhost:8001/pipeline/ingest -H "Content-Type: application/json" -d '{}'

# 5. Frontend (terminal 3)
cd frontend-client && npm install && npm run dev
```

Open http://localhost:5173 and log in.
