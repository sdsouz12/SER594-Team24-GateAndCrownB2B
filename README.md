# Gate & Crown B2B Platform

Repository:
https://github.com/sdsouz12/SER594-Team24-GateAndCrownB2B

Gate & Crown is a full-stack B2B ordering platform built using Go (Gin) and Vue 3, providing secure authentication, PostgreSQL persistence, and REST-based communication between frontend and backend services.


## Tech Stack

### Backend
- Go 1.25+
- Gin (HTTP framework)
- PostgreSQL (pgx driver)
- REST / JSON API
- golang-migrate (SQL migrations)
- JWT authentication
- bcrypt password hashing

### Frontend
- Vue 3
- Vue Router
- Vite
- Tailwind CSS v4
- Native Fetch API


## Requirements

Install before running the project:

### Backend
- Go 1.25+
- PostgreSQL
- Database user with CREATEDB permission

### Frontend
- Node.js LTS
- npm


## Backend Setup

### 1. Configure Environment Variables

Create:

backend-service/.env

Example:

DATABASE_URL=postgres://username:password@localhost:5432/gatecrown?sslmode=disable
JWT_SECRET=super-secret-key
FRONTEND_URL=http://localhost:5173

Template:

backend-service/.env.example


### 2. Install Dependencies

Run from repository root:

go -C backend-service mod tidy


### 3. Run Database Migrations

go run ./backend-service/cmd/migrate up


### 4. Start Backend Server

go run ./backend-service/cmd/server

Backend runs at:

http://localhost:8080

Important:

Always start backend using:

go run ./backend-service/cmd/server

Do NOT run main.go directly.


## Frontend Setup

Start backend first.

Then run:

cd frontend-client

npm install

npm run dev

Frontend runs at:

http://localhost:5173

Development server proxies:

/api → http://localhost:8080

Environment template:

frontend-client/.env.example


## Project Structure

SER594-Team24-GateAndCrownB2B
│
├── backend-service
│   ├── cmd/server
│   ├── cmd/migrate
│   ├── internal
│   └── migrations
│
├── frontend-client
│   ├── src
│   └── public
│
└── go.work


## Authentication Flow

User Login
↓
JWT Issued (Backend)
↓
Stored in Client
↓
Attached to API Requests
↓
Protected Route Access

Security features include:

- bcrypt password hashing
- JWT authentication middleware
- protected REST endpoints


## Development Workflow

Start backend:

go run ./backend-service/cmd/server

Start frontend:

npm run dev

Open application:

http://localhost:5173


## Recommended Future Improvements (Optional)

To strengthen this project for internship or portfolio presentation:

- Docker support
- Swagger/OpenAPI documentation
- Refresh-token authentication flow
- Role-based authorization
- Structured logging
- Unit tests for handlers and services
- GitHub Actions CI pipeline
