# Personal Web Platform

A personal website combining a professional portfolio (CV) and an AI-focused blog.

## Features

- **Public Section:**
  - Author profile (CV) with contacts and bio.
  - Responsive design (Mobile/Desktop).
  - Light and Dark themes.
- **Blog:**
  - Article feed with Markdown and HTML support.
  - Read time estimation.
  - Nested comments and likes.
  - Infinite scroll.
- **Admin Panel:**
  - Profile management.
  - Post CRUD with WYSIWYG editor.
  - Image upload.
- **Security:**
  - OAuth authentication (Google, GitHub, VK ID).
  - Role-Based Access Control (Admin/User).
  - Rate Limiting.

## Tech Stack

### Backend

- **Language:** Go 1.25
- **Framework:** Chi v5
- **Database:** PostgreSQL 15
- **Auth:** OAuth 2.0/2.1

### Frontend

- **Framework:** React 18, Vite
- **Language:** TypeScript
- **UI:** Tailwind CSS, shadcn/ui
- **State:** TanStack Query

### Infrastructure

- **Containerization:** Docker, Docker Compose
- **Web Server:** Nginx
- **CI/CD:** GitHub Actions

## Quick Start

**Requirements:** Docker, Docker Compose, Make.

1. **Initialization**

   ```bash
   make init
   ```

2. **Run Development Environment**
   ```bash
   make dev
   ```
   - Website: http://localhost
   - Backend API: http://localhost:8080
   - Frontend: http://localhost:5173

## Makefile Commands

| Command      | Description                                    |
| ------------ | ---------------------------------------------- |
| `make init`  | Initial setup                                  |
| `make dev`   | Start development environment                  |
| `make logs`  | View logs                                      |
| `make stop`  | Stop containers                                |
| `make reset` | Reset environment (delete containers and data) |
| `make check` | Run linters                                    |
| `make test`  | Run backend tests                              |

## Access Management

Default users have the `user` role. To assign admin rights:

```bash
docker compose exec -T db psql -U postgres -d pwp_db -c "UPDATE users SET role = 'admin' WHERE email = 'YOUR_EMAIL@gmail.com'"
```

## Project Structure

```
.
├── backend/        # Go API server
│   ├── config/     # Configuration
│   └── migrations/ # SQL migrations
├── frontend/       # React application
├── docs/           # Documentation
├── docker-compose.yml
├── nginx.conf
└── Makefile
```

## Deployment

The project supports two deployment modes: Local and Production.

### 1. Configuration

Create a .env file from the example:

```bash
cp .env.example .env
```

Fill in the sensitive data in .env.

### 2. Local Development

Uses docker-compose.yml. Builds with hot-reload (Air for Go, Vite for React).

```bash
# Start local dev
make dev
```

App available at: http://localhost

### 3. Production Deployment

Uses docker-compose.prod.yml. Optimized builds, Nginx with SSL (Let's Encrypt).

Server Requirements:

- Docker & Docker Compose installed
- Ports 80 and 443 open
- Domain DNS pointed to server IP (5.53.125.146 -> fus1ond.ru)

Steps:

1. Clone & Config:

   ```bash
   git clone <repo_url>
   cd curriculum_vitae
   cp .env.example .env
   cp backend/config/production.yaml.example backend/config/production.yaml
   ```

2. Edit .env for Production:

   ```ini
   ENV=production
   DOMAIN=fus1ond.ru
   ACME_EMAIL=your@email.com
   OAUTH_BASE_URL=https://fus1ond.ru
   OAUTH_FRONTEND_URL=https://fus1ond.ru
   COOKIE_DOMAIN=fus1ond.ru
   ... set other secrets (DB, OAuth)
   ```

3. Run with SSL Setup:

   For the first run (to generate certificates), use the helper script:

   ```bash
   chmod +x init-letsencrypt.sh
   ./init-letsencrypt.sh
   ```

   For subsequent updates/restarts:

   ```bash
   docker compose -f docker-compose.prod.yml up -d --build
   ```
