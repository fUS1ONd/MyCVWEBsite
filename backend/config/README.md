# Backend Configuration

All backend configuration is stored in YAML files in this directory.

## Configuration Files

### Local Development
- **local.yaml.example** - Template for local development
- **local.yaml** - Your actual local configuration (git-ignored)

### Production
- **production.yaml.example** - Template for production
- **production.yaml** - Your actual production configuration (git-ignored)

## Quick Start

### Local Development

1. Copy the example file:
   ```bash
   cp backend/config/local.yaml.example backend/config/local.yaml
   ```

2. Edit `local.yaml` with your local settings:
   - Database credentials
   - OAuth provider credentials (if testing auth)
   - Session secret

3. The backend will automatically use `config/local.yaml` when running with docker-compose

### Production Deployment

1. Copy the example file:
   ```bash
   cp backend/config/production.yaml.example backend/config/production.yaml
   ```

2. Edit `production.yaml` with your production settings:
   - **Database URL** with strong password and SSL enabled
   - **Session secret** - generate with `openssl rand -base64 48`
   - **OAuth credentials** for Google, GitHub, VK
   - **Cookie domain** - your production domain (e.g., `fus1ond.ru`)
   - **Base URLs** - your HTTPS domain

3. Ensure `production.yaml` is properly secured:
   ```bash
   chmod 600 backend/config/production.yaml
   ```

## Configuration Structure

All backend settings are in YAML files:
- Database connection
- HTTP server settings
- Authentication & session management
- OAuth providers (Google, GitHub, VK)
- CORS settings
- Rate limiting
- Profile information

## Environment Variables

`.env` file is now minimal and only contains:
- **Docker-compose variables**: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- **Frontend variables**: `VITE_BACKEND_URL`, `VITE_GTM_ID`
- **Infrastructure**: `DOMAIN`, `ACME_EMAIL`

The backend itself reads ALL configuration from the YAML file specified by `CONFIG_PATH` environment variable.

## How It Works

1. Docker-compose sets `CONFIG_PATH` environment variable:
   - Local: `CONFIG_PATH=config/local.yaml`
   - Production: `CONFIG_PATH=config/production.yaml`

2. Backend reads the specified YAML file at startup

3. All settings (including database URL, OAuth secrets, etc.) come from the YAML file

## Security Notes

⚠️ **IMPORTANT**:
- Never commit `local.yaml` or `production.yaml` to git
- Keep `production.yaml` secure (chmod 600)
- Use strong random values for `session_secret`
- Enable SSL (`sslmode=require`) for production database
- Set `cookie_secure: true` in production

## Migration from Old Structure

Old approach (environment variables everywhere):
```yaml
environment:
  - DATABASE_URL=${DATABASE_URL}
  - SESSION_SECRET=${SESSION_SECRET}
  - OAUTH_BASE_URL=${OAUTH_BASE_URL}
  # ... many more variables
```

New approach (single config file):
```yaml
environment:
  - CONFIG_PATH=config/production.yaml
```

All settings are now in the YAML file for better organization and easier management.
