# Yojana Portal Documentation

## Overview

**Yojana Portal** is a full‑stack web application that helps users discover and apply for government jobs. The project consists of:

- **Backend** (Go) – REST API handling authentication, job listings, eligibility checks, and admin functionality.
- **Frontend** (Vue 3 + Vite) – Single‑page application that consumes the API and provides a modern UI.

Both parts are developed together in this repository. This documentation explains how to set up the development environment, run the services locally, and perform common tasks.

---

## Prerequisites

| Tool | Minimum version |
|------|-----------------|
| **Go** | 1.22 |
| **Node.js** | 20.x |
| **npm** | 10.x |
| **git** | any recent version |
| **Docker** (optional) | 27.x |

> The project is primarily targeted at Windows, but the commands work on macOS/Linux as well.

---

## Repository Structure

```
.
├─ backend/                 # Go server source code
│   ├─ cmd/server/main.go   # Entry point of the API server
│   ├─ internal/            # Handlers, middleware, DB layer, etc.
│   └─ go.mod / go.sum      # Go module files
├─ frontend/                # Vue 3 + Vite SPA
│   ├─ src/                 # Application source code
│   ├─ public/              # Static assets (favicon, icons)
│   ├─ vite.config.js       # Vite configuration
│   └─ package.json         # npm dependencies
├─ Documentations/          # Project documentation (this folder)
│   ├─ README.md            # **You are here**
│   └─ api_spec.md          # OpenAPI specification for the backend API
├─ .env.example             # Example environment variables
└─ README.md                # Top‑level project readme (frontend entry point)
```

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-org/yojana-portal.git
cd yojana-portal
```

### 2. Set up environment variables

Copy the example file and adjust values as needed:

```bash
cp .env.example .env
```

Key variables:
- `PORT` – Port for the backend server (default **8080**).
- `DB_DSN` – Data source name for the PostgreSQL database.
- `JWT_SECRET` – Secret used to sign JWT tokens.
- `ADMIN_USERS` – Comma‑separated list of admin email addresses.

### 3. Run the backend

```bash
# Navigate to the backend directory
cd backend

# Install Go dependencies (handled automatically by go build/run)
# Start the server
go run ./cmd/server
```

The API will be available at `http://localhost:8080`.

### 4. Run the frontend

Open a new terminal window/tab:

```bash
# From the repository root
cd frontend

# Install npm dependencies
npm install

# Start the dev server
npm run dev
```

The SPA will be served at `http://localhost:5173` (or the port shown in the console). The UI expects the backend to be reachable at `http://localhost:8080`; you can change the proxy target in `vite.config.js` if needed.

---

## Database Setup (PostgreSQL)

The backend uses PostgreSQL via the `pgx` driver. A quick local setup using Docker:

```bash
docker run -d \
  -e POSTGRES_USER=yojana \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=yojana \
  -p 5432:5432 \
  --name yojana-postgres \
  postgres:16-alpine
```

Update `DB_DSN` in `.env` to something like:
```
postgres://yojana:secret@localhost:5432/yojana?sslmode=disable
```

Run the DB migrations (if any) – currently the project creates tables on first start, but you can also execute the SQL scripts under `backend/internal/db/migrations/`.

---

## Testing

### Backend tests

```bash
cd backend
go test ./...   # runs all unit tests
```

### Frontend tests

```bash
cd frontend
npm run test   # runs Jest/Vitest suite
```

---

## Docker (optional all‑in‑one)

A `Dockerfile` exists for the backend. You can build and run both services with Docker Compose (not included by default, but an example is provided below):

```yaml
# docker-compose.yml (example)
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: yojana
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: yojana
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    env_file: .env
    depends_on:
      - db
    ports:
      - "8080:8080"

  frontend:
    build: ./frontend
    ports:
      - "5173:5173"
    depends_on:
      - backend
```

Run with:
```bash
docker compose up --build
```

---

## API Specification

The OpenAPI spec is located in `Documentations/api_spec.md`. It describes all endpoints, request/response schemas, authentication requirements, and example payloads. Tools such as **Swagger UI** or **Redoc** can render this file for interactive exploration.

---

## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your‑feature`).
3. Make your changes and ensure both backend and frontend tests pass.
4. Submit a Pull Request against the `main` branch.
5. Follow the code‑style guidelines (gofmt for Go, eslint/prettier for JavaScript/TypeScript).

---

## License

This project is licensed under the **MIT License** – see the `LICENSE` file at the repository root.

---

## Contact & Support

For questions, open an issue on GitHub or reach out to the maintainers via the `#yojana-portal` channel on the organization’s Slack workspace.
