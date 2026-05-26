# Yojana Portal

Yojana Portal is a full-stack government scheme discovery application. It helps users browse public schemes, check starter eligibility matches, save or track schemes, and use an assistant-style chat surface for scheme questions.

The project is split into a Go API backend and a Vue frontend.

## Tech Stack

| Area | Technology |
| --- | --- |
| Backend | Go 1.22, Gin, GORM |
| Database | PostgreSQL, Supabase-compatible connection URLs |
| Auth | JWT bearer tokens, bcrypt password hashing |
| Email | SendGrid utility for OTP delivery |
| Frontend | Vue 3, Vite, Vue Router, Pinia, Axios |
| Styling | Tailwind CSS plus project CSS |

## Project Structure

```text
yojana-portal/
  backend/
    config/          App configuration and database connection
    handlers/        HTTP request handlers
    middleware/      Auth and CORS middleware
    models/          GORM models and API data shapes
    routes/          API route registration
    utils/           JWT, OTP, email, and response helpers
    main.go          Backend entrypoint
  frontend/
    src/
      assets/        Global styles and static assets
      components/    Shared UI and feature components
      locales/       English, Hindi, and Marathi locale files
      router/        Vue Router routes
      services/      Axios API clients
      stores/        Pinia stores
      views/         Page-level Vue views
    vite.config.js
  README.md
```

## Features

- Public scheme browsing and filtering
- Scheme detail API endpoint
- User registration and login with JWT sessions
- Password hashing with bcrypt
- OTP creation with hashed storage
- Auth-protected user profile, saved schemes, tracked schemes, chatbot, and notification routes
- Admin-protected scheme and user management routes
- CORS configuration for frontend/backend development
- Graceful backend shutdown and HTTP timeouts
- Multilingual frontend locale files for English, Hindi, and Marathi

## Current Status

The backend exposes real Gin routes and database-backed handlers. Some frontend screens still use starter local data while the service files are ready for API integration. The chatbot endpoint is credential-gated and currently returns a configured/not-configured response until the AI provider integration is completed.

## Prerequisites

- Go 1.22 or newer
- Node.js 18 or newer
- npm
- PostgreSQL database URL, such as a Supabase connection string

## Backend Setup

Create `backend/.env`:

```env
APP_ENV=development
APP_PORT=8080
DB_URL=postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require
AUTO_MIGRATE=true

JWT_SECRET=replace-with-a-long-random-secret
JWT_EXPIRY_HOURS=24

FRONTEND_URL=http://localhost:5173
CORS_ALLOWED_ORIGINS=http://localhost:5173

SENDGRID_API_KEY=
FROM_EMAIL=
FROM_NAME=Yojana Portal

ANTHROPIC_API_KEY=
```

Notes:

- `DB_URL` can also be provided as `DATABASE_URL`.
- `APP_PORT` can also be provided as `PORT`.
- `FRONTEND_URL` can also be provided as `FRONTEND_ORIGIN`.
- `ANTHROPIC_API_KEY` can also be provided as `CLAUDE_API_KEY`.
- `AUTO_MIGRATE` defaults to `true` in development and `false` in production.
- In production, use a strong `JWT_SECRET` with at least 32 characters.

Run the backend:

```powershell
cd backend
go mod tidy
go run .
```

Health check:

```powershell
Invoke-WebRequest http://localhost:8080/health
```

## Frontend Setup

Optional frontend environment file:

```env
VITE_API_URL=http://localhost:8080/api
```

Run the frontend:

```powershell
cd frontend
npm install
npm run dev
```

The Vite app runs at:

```text
http://localhost:5173
```

## Useful Commands

Backend:

```powershell
cd backend
go test ./...
go vet ./...
gofmt -w .
```

Frontend:

```powershell
cd frontend
npm run build
npm run preview
```

## API Overview

All API endpoints are served from `/api` unless noted otherwise.

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| GET | `/health` | Public | Service health check |
| POST | `/api/auth/register` | Public | Create user account |
| POST | `/api/auth/login` | Public | Login and receive JWT |
| POST | `/api/auth/otp` | Public | Create OTP for a contact |
| POST | `/api/auth/logout` | Authenticated | Logout response endpoint |
| GET | `/api/schemes` | Public | List published schemes |
| GET | `/api/schemes/search` | Public | Search/filter published schemes |
| GET | `/api/schemes/:id` | Public | Fetch one published scheme |
| GET | `/api/categories` | Public | List scheme categories |
| POST | `/api/eligibility/check` | Public | Check starter scheme matches |
| GET | `/api/user/profile` | Authenticated | Fetch current user's profile |
| GET | `/api/user/saved-schemes` | Authenticated | Fetch saved schemes |
| GET | `/api/user/tracked-schemes` | Authenticated | Fetch tracked applications |
| POST | `/api/chatbot` | Authenticated | Chatbot request endpoint |
| GET | `/api/notifications` | Authenticated | List current user's notifications |
| POST | `/api/notifications/deadline-reminder` | Authenticated | Queue deadline reminder |
| GET | `/api/admin` | Admin | Admin dashboard counts |
| POST | `/api/admin/schemes` | Admin | Create scheme |
| PUT | `/api/admin/schemes/:id` | Admin | Update scheme |
| DELETE | `/api/admin/schemes/:id` | Admin | Delete scheme |
| GET | `/api/admin/users` | Admin | List users |

Authenticated requests must include:

```text
Authorization: Bearer <token>
```

## Response Format

Successful responses use:

```json
{
  "success": true,
  "message": "operation completed",
  "data": {}
}
```

Error responses use:

```json
{
  "success": false,
  "message": "error message"
}
```

## Development Notes

- Keep real credentials in `backend/.env`. It is ignored by Git.
- Do not return OTP codes from the API in production flows. The backend stores hashed OTP values and sends email only when SendGrid is configured.
- Public scheme routes currently return records with `status = "published"`.
- Admin-only routes require a JWT whose role is `admin`.
- Database migrations are handled through GORM `AutoMigrate` when `AUTO_MIGRATE=true`.

## Production Checklist

- Set `APP_ENV=production`.
- Use a strong `JWT_SECRET`.
- Set `AUTO_MIGRATE=false` and run reviewed migrations separately.
- Restrict `CORS_ALLOWED_ORIGINS` to deployed frontend domains.
- Configure SendGrid sender identity before enabling OTP email delivery.
- Add rate limiting for login and OTP routes.
- Add automated tests for auth, admin authorization, and scheme CRUD.
