# Gate & Crown B2B Platform

SER 594 — AI for Software Engineers | Final Milestone

Full-stack B2B ordering platform with semantic search, LLM order assistant, and RAG — built with Go, Vue 3, and Python (FastAPI).

---

## Deployed / Running Locally

This project ships with a complete Docker Compose setup. See **Docker Setup** below for the one-command run.

---

## Requirements

- Go 1.25+
- Node.js LTS + npm
- PostgreSQL 16+ with pgvector extension
- Python 3.11 (via Anaconda/conda recommended)
- Docker + Docker Compose (for containerised run)

---

## Option A — Docker (Recommended)

No manual setup required — Docker handles the database, migrations, backend, AI service, and frontend.

```bash
docker compose up --build
```

Then open **http://localhost:3000** and log in.

> First run downloads the embedding model (~66 MB) inside the AI-service container — allow ~2 minutes.

To stop:

```bash
docker compose down
```

To wipe data and start fresh:

```bash
docker compose down -v
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

Install dependencies and start:

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
3. After login you land on the **Demo Dashboard** which shows:
   - Data pipeline status
   - Semantic vector search
   - LLM Order Assistant (Claude)
   - RAG Q&A
4. Navigate to **/catalog** to browse products and use AI-powered search

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

### Backend (Go) — 15 unit tests

```bash
cd backend-service
go test ./... -v
```

### AI Service (Python) — 15 unit tests

```bash
cd AI-service
pip install -r requirements-test.txt
pytest tests/ -v
```

All 30 tests run without a database or network connection (external I/O is mocked).

CI runs both suites automatically on every push via GitHub Actions (`.github/workflows/ci.yml`).

---

## AI Techniques Implemented

| Technique | Description | Endpoint |
|---|---|---|
| Vector Search | BAAI/bge-small-en-v1.5 embeddings + pgvector cosine similarity | `POST /api/catalog/search` |
| LLM Assistant | Claude Haiku returns structured JSON order recommendations | `POST /api/ai/assist` |
| RAG | Top-5 relevant products retrieved as context; Claude answers grounded Q&A | `POST /api/ai/rag` |

---

## Project Structure

```
gate-crown/
├── backend-service/          Go + Gin REST API
│   ├── cmd/server/           Server entry point
│   ├── cmd/migrate/          Database migration tool
│   ├── internal/auth/        Auth (login, register, JWT) + unit tests
│   ├── internal/catalog/     Catalog (list + semantic search)
│   ├── internal/aiproxy/     Proxy to Python AI service
│   ├── migrations/           SQL migration files (001–006)
│   ├── Dockerfile
│   └── entrypoint.sh
├── frontend-client/          Vue 3 + Vite + Tailwind CSS
│   ├── src/views/
│   │   ├── WelcomeView.vue   Demo dashboard
│   │   ├── CatalogView.vue   Product catalog + AI search
│   │   ├── LoginView.vue
│   │   └── RegisterView.vue
│   ├── Dockerfile
│   └── nginx.conf
├── AI-service/               Python FastAPI AI microservice
│   ├── main.py               Vector search, LLM assistant, RAG
│   ├── tests/test_main.py    15 unit tests
│   ├── requirements.txt
│   ├── requirements-test.txt
│   └── Dockerfile
├── docker-compose.yml
├── .github/workflows/ci.yml  GitHub Actions CI
└── go.work
```

---

## Quick Start (manual, all commands in order)

```bash
# 1. Migrations
go run ./backend-service/cmd/migrate up

# 2. Backend (terminal 1)
go run ./backend-service/cmd/server

# 3. AI service (terminal 2)
conda activate gatecrown
/opt/anaconda3/envs/gatecrown/bin/python AI-service/main.py

# 4. Ingest embeddings (once)
curl -X POST http://localhost:8001/pipeline/ingest -H "Content-Type: application/json" -d '{}'

# 5. Frontend (terminal 3)
cd frontend-client && npm install && npm run dev
```

Open http://localhost:5173 and log in.
