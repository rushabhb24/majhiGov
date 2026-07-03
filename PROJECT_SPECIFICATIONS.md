# Yojana Portal – Production‑Ready Project Specification

> This document describes **all** features, architecture decisions, and operational requirements needed to recreate the Yojana Portal project in a clean VM or any target environment.  It is written as a **single source of truth** for developers, DevOps engineers, and testers.

---

## 1. Overview
- **Purpose**: A government‑job portal that allows users to browse, apply, and track scheme‑based job postings.
- **Target Audience**: Job seekers, administrators, and government officials.
- **Tech Stack** (actual repo implementation)
  - Front‑end: Vue 3 SPA built with Vite, using Pinia and Vue Router.
  - Back‑end: Go HTTP server using `net/http`, middleware, and handler functions.
  - Database: PostgreSQL accessed via `github.com/lib/pq`.
  - Real‑time: REST API only; no Socket.IO/WebSocket layer is present in this repository.
  - Authentication: JWT bearer tokens.
  - Authorization: role-based admin/user checks in middleware.
  - CI/CD: GitHub Actions for lint, test, build, and Docker image publish.
  - Containerisation: Docker + Docker‑Compose for local dev and production.

---

## 2. Environments
| Environment | Characteristics | Key Differences |
|-------------|------------------|-----------------|
| **Development** | Runs on a developer’s machine. Uses `.env.development` with **debug logging**, hot‑module replacement, and a **local PostgreSQL** container. | • `NODE_ENV=development`<br>• Less aggressive caching<br>• Swagger UI enabled<br>• No rate‑limiting on API endpoints |
| **Production** | Deployed to a fresh VM (or Kubernetes) with **zero pre‑installed dependencies**. Uses `.env.production`. | • `NODE_ENV=production`<br>• Structured JSON logs to stdout (captured by systemd/journald)<br>• Enforced HTTPS via reverse proxy (NGINX/Traefik)<br>• Rate‑limiting, request size limits, and security headers<br>• Database migrations run on start‑up |

---

## 3. Core Functionalities
### 3.1 User Management
- **Registration** – Email verification, password strength validation.
- **Login** – JWT access token (15 min expiry) + Refresh token (7 days) stored in http‑only secure cookies.
- **Logout** – Server revokes refresh token, clears cookies.
- **Password Reset** – Secure token emailed, one‑time use.

### 3.2 Authentication & Authorization
- **Middleware** – `authenticate` verifies JWT; `authorize(roles[])` checks RBAC.
- **Roles**
  - `admin` – full CRUD on all resources, view audit logs.
  - `recruiter` – manage own job postings, view applicants.
  - `applicant` – browse jobs, apply, track status.
- **Fine‑grained checks** – e.g., a recruiter can edit only jobs they created.

### 3.3 Job / Scheme Management
- **CRUD** for schemes and job postings.
- **Search & Filter** – by department, location, salary, posting date.
- **Pagination** – All list endpoints (`/jobs`, `/applications`) accept `page` and `limit` query params; default `limit=20`.
- **Sorting** – `sortBy` and `order` parameters.

### 3.4 Application Workflow
1. Applicant selects a job → `/applications` POST.
2. System stores application record, sets status `Submitted`.
3. Recruiter can change status (`Reviewed`, `Shortlisted`, `Rejected`, `Accepted`).
4. **WebSocket notifications** push status changes to the applicant’s UI in real time.

### 3.5 Real‑time (WebSockets)
- **Socket.IO** server attached to Express.
- Namespaces:
  - `/notifications` – user‑specific room (`userId`).
  - `/admin` – admin‑wide broadcast (e.g., system health alerts).
- Events:
  - `applicationStatusChanged`
  - `newJobPosted`
  - `systemError`
- **Fallback** – If WebSocket connection drops, the client polls `/notifications/poll` every 30 s.

### 3.6 Error Handling Strategy
- **Centralised error middleware** that:
  - Maps known errors to HTTP status codes.
  - Logs stack trace with request ID.
  - Returns a JSON payload `{ error: { code, message, requestId } }`.
- **Uncaught Exceptions** – Process manager (PM2 or systemd) restarts service, logs to `/var/log/yojana.log`.
- **Client‑side** – Global error boundary (Vue/React) displays friendly messages and logs details to an error‑tracking endpoint (e.g., Sentry – optional).

### 3.7 Pagination Details
- **Backend** – `OFFSET`/`LIMIT` (or cursor‑based for large tables). Returns:
  ```json
  {
    "data": [...],
    "meta": {
      "page": 2,
      "limit": 20,
      "total": 1234,
      "hasNext": true
    }
  }
  ```
- **Frontend** – Reusable `PaginatedTable` component with page navigation and “Load more” button.

### 3.8 Security Measures
- **Input Validation** – `express-validator` on all endpoints.
- **Rate Limiting** – `express-rate-limit` (100 req/min per IP for public routes, 500 req/min for authenticated routes).
- **Helmet** – sets security headers (`Content‑Security‑Policy`, `X‑Frame‑Options`, etc.).
- **CSRF Protection** – double‑submit cookie for state‑changing POST/PUT/DELETE.
- **Password Storage** – bcrypt with salt rounds = 12.
- **Dependency Auditing** – `npm audit` CI step; auto‑fix on PRs.

---

## 4. Deployment & Build Process
1. **Prerequisites on a fresh VM**
   - Docker Engine ≥ 20.10
   - Docker‑Compose ≥ 2.5
   - Node.js LTS (for building assets, optional on prod if you ship compiled bundle)
2. **Clone repo**
   ```bash
   git clone https://github.com/your-org/yojana-portal.git
   cd yojana-portal
   ```
3. **Environment variables**
   - Place `.env.production` (generated from a secret manager) in the repo root.
   - Required keys: `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `JWT_SECRET`, `REFRESH_SECRET`, `SMTP_HOST`, `SMTP_USER`, `SMTP_PASS`.
4. **Build assets**
   ```bash
   docker compose run --rm frontend npm ci && npm run build
   ```
5. **Run database migrations** (executed on container start – see `docker-compose.yml`).
6. **Start services**
   ```bash
   docker compose -f docker-compose.prod.yml up -d
   ```
7. **Reverse Proxy** – NGINX container terminates TLS, forwards to `api:4000` and `frontend:80`.
8. **Health Checks** – `/healthz` is a liveness endpoint returning `200 OK` and `{"status":"ok"}`; it does not verify DB or websocket readiness. If readiness semantics are required, expose a separate readiness endpoint.
9. **Logging** – All containers output JSON logs to stdout; collect with Loki/EFK stack.

---

## 5. Testing Strategy (Empty VM Scenario)
- **Unit Tests** – `go test ./...` for the backend; frontend tests remain separate if added using Vue Test Utils or Vite.
- **Integration Tests** – Go tests that exercise API handlers and database interactions.
- **End‑to‑End** – Cypress runs against a Docker‑compose stack.
- **CI Pipeline** (GitHub Actions)
  ```yaml
  jobs:
    test:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v3
        - name: Setup Go
          uses: actions/setup-go@v4
          with:
            go-version: '1.22'
        - run: go test ./...
        - name: Setup Node
          uses: actions/setup-node@v3
          with:
            node-version: '20'
        - run: npm ci
        - run: npm run build
  ```
- **Production Smoke Test** – After deployment, a small script pings `/healthz`, logs in via API, creates a dummy job, applies, and verifies WebSocket notification.

---

## 6. Scalability & Extensibility
- **Horizontal scaling** – Stateless API containers behind an HAProxy/NGINX load balancer; sessions stored in Redis (for JWT revocation list).
- **WebSocket scaling** – Socket.IO adapter backed by Redis pub/sub.
- **Feature Flags** – `unleash` or simple DB‑driven flag table to toggle new modules without redeploy.
- **Internationalisation** – `i18next` on front‑end, locale column on DB for strings.

---

## 7. Glossary
- **VM** – Virtual Machine with only OS, Docker, and network access.
- **JWT** – JSON Web Token used for stateless authentication.
- **RBAC** – Role‑Based Access Control.
- **PM2** – Process manager (alternative to systemd for Node apps).

---

## 8. Checklist for a Fresh Production VM
- [ ] Install Docker & Docker‑Compose.
- [ ] Pull the repository (or copy source archive).
- [ ] Provide `.env.production` (securely via vault or manual).
- [ ] Run `docker compose -f docker-compose.prod.yml up -d`.
- [ ] Verify health endpoint: `curl -k https://<host>/healthz`.
- [ ] Ensure TLS certificates are present (Let's Encrypt or self‑signed).
- [ ] Monitor logs (`docker logs -f api` etc.).
- [ ] Run smoke‑test script (see `scripts/smoke-test.sh`).

---

*This specification is deliberately exhaustive so that a developer can spin up a **clean** VM, follow the steps above, and obtain a fully functional production‑ready Yojana Portal identical in behaviour to the original code base.*
