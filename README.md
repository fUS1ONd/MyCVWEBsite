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

## Getting Started

Prerequisites: Docker, Docker Compose, Make.

### Local Development

1. **Initialize the project:**

   ```bash
   make init
   ```
   This sets up the necessary configuration files.

2. **Start development environment:**

   ```bash
   make dev
   ```

   The application will be available at:
   - **Frontend:** http://localhost:5173
   - **Backend API:** http://localhost:8080

### Production Deployment

To deploy in production mode with Nginx and SSL support:

```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

Ensure you have configured `backend/config/production.yaml` and set the necessary environment variables in `.env` (domain, database credentials, etc.).

## Useful Commands

- `make logs` - View logs
- `make stop` - Stop containers
- `make reset` - Reset environment (delete containers and volumes)

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
