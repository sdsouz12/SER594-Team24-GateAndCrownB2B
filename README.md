# Gate & Crown B2B Platform

SER 594 — AI for Software Engineers | Milestone 2

Full-stack B2B ordering platform with semantic search, LLM order assistant, and RAG — built with Go, Vue 3, and Python (FastAPI).

---

## Requirements

Before starting, make sure you have:

- Go 1.25+
- Node.js LTS + npm
- PostgreSQL (with pgvector extension)
- Python 3.11 (via Anaconda/conda recommended)
- Conda (Anaconda)

---

## Step 1 — Install pgvector

pgvector is required for AI semantic search.

```bash
brew install pgvector
```

Then enable it in PostgreSQL:

```bash
psql postgres -c "CREATE EXTENSION IF NOT EXISTS vector;"
```

---

## Step 2 — Configure Backend Environment

Copy and edit the backend environment file:

```bash
cp backend-service/.env.example backend-service/.env
```

Open `backend-service/.env` and set your database credentials:
```bash
open backend-service/.env
```

```
DATABASE_URL=postgresql://localhost:5432/gate_crown
PORT=8080
JWT_SECRET=change-me-to-a-long-random-secret
JWT_EXPIRATION_HOURS=24
FRONTEND_URL=http://localhost:5173
AI_SERVICE_URL=http://localhost:8001
```

> Change `DATABASE_URL` to match your PostgreSQL username if needed, e.g. `postgresql://youruser@localhost:5432/gate_crown`

---

## Step 3 — Run Database Migrations

This creates the database and runs all migrations (auth tables, organizations, catalog products with pgvector):

```bash
go run ./backend-service/cmd/migrate up
```

---

## Step 4 — Start the Backend

```bash
go run ./backend-service/cmd/server
```

Backend runs at: http://localhost:8080

---

## Step 5 — Set Up the AI Service

The AI service handles vector search, LLM assistant, and RAG using Python.

### Create conda environment

Make sure you have anaconda installed:

After this make sure you accept the terms and conditions by running the commands that show up on your terminal 

Init conda for your shell:

```bash
conda init zsh
```
Ensure that you restart your terminal after this step

```bash
conda create -n gatecrown python=3.11 -y
conda activate gatecrown
```

### Configure AI service environment

```bash
cp AI-service/.env.example AI-service/.env
```

The `.env.example` already contains the Anthropic API key — no changes needed.

### Install dependencies

```bash
cd AI-service
/opt/anaconda3/envs/gatecrown/bin/pip install -r requirements.txt
```

### Start the AI service

```bash
/opt/anaconda3/envs/gatecrown/bin/python main.py
```

AI service runs at: http://localhost:8001

---

## Step 6 — Ingest Catalog Embeddings (In a new terminal)

This generates vector embeddings for all 20 catalog products and stores them in PostgreSQL. Run once:

```bash
curl -X POST http://localhost:8001/pipeline/ingest \
  -H "Content-Type: application/json" \
  -d '{}'
```

Expected response:

```json
{"message": "Ingested 20 products.", "processed": 20}
```

---

## Step 7 — Start the Frontend

```bash
cd frontend-client
npm install
npm run dev
```

Frontend runs at: http://localhost:5173

---

## Step 8 — Use the Application

1. Open http://localhost:5173
2. Click **Register** to create a new account, or use the test account:
   - Username: `admin`
   - Password: `admin123`
3. After login you are taken to the **Milestone 2 Demo Dashboard** which shows:
   - Data pipeline status (products embedded)
   - Semantic vector search — type any natural language query
   - LLM Order Assistant — Claude recommends products based on your description
   - RAG — Claude answers questions using retrieved catalog context
4. Visit **/catalog** from the top navigation to browse all products and use AI search

---

## Test Accounts

| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | SUPERADMIN |
| john.doe | admin123 | ADMIN (Acme Corp) |
| jane.smith | admin123 | ADMIN (Acme Corp) |
| sarah.jones | admin123 | ADMIN (Beta Solutions) |

---

## AI Techniques Implemented

| Technique | Description | Endpoint |
|-----------|-------------|----------|
| Vector Search | Semantic search using BAAI/bge-small-en-v1.5 embeddings + pgvector cosine similarity | `POST /api/catalog/search` |
| LLM Assistant | Claude Haiku returns structured JSON order configuration | `POST /api/ai/assist` |
| RAG | Retrieves top-5 relevant products as context, Claude generates grounded answer | `POST /api/ai/rag` |

---

## Project Structure

```
gate-crown/
├── backend-service/       Go + Gin REST API
│   ├── cmd/server/        Server entry point
│   ├── cmd/migrate/       Database migration tool
│   ├── internal/auth/     Authentication (login, register, JWT)
│   ├── internal/catalog/  Catalog module (list + search)
│   ├── internal/aiproxy/  Proxy to Python AI service
│   └── migrations/        SQL migration files
├── frontend-client/       Vue 3 + Vite + Tailwind CSS
│   └── src/views/
│       ├── LoginView.vue
│       ├── RegisterView.vue
│       ├── CatalogView.vue
│       └── WelcomeView.vue  ← Milestone 2 demo dashboard
├── AI-service/            Python FastAPI AI service
│   └── main.py            Vector search, LLM assistant, RAG
└── go.work
```

---

## Quick Start (all steps in order)

```bash
# 1. Migrations
go run ./backend-service/cmd/migrate up

# 2. Backend
go run ./backend-service/cmd/server

# 3. AI service (new terminal)
conda activate gatecrown
cd AI-service
/opt/anaconda3/envs/gatecrown/bin/python main.py

# 4. Ingest embeddings (once)
curl -X POST http://localhost:8001/pipeline/ingest -H "Content-Type: application/json" -d '{}'

# 5. Frontend (new terminal)
cd frontend-client
npm install && npm run dev
```

Open http://localhost:5173 and log in.
