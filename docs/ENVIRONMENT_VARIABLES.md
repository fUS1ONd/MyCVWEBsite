# Environment Variables Reference

Полный справочник всех environment variables, используемых приложением.

## Основные настройки

### Backend Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `CONFIG_PATH` | No | `config/local.yaml` | Путь к файлу конфигурации |
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |

**Пример:**
```bash
CONFIG_PATH=config/production.yaml
DATABASE_URL=postgres://user:password@db-host:5432/dbname?sslmode=require
```

---

## Authentication & Security

### Session Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SESSION_SECRET` | Yes | - | Секретный ключ для шифрования сессий (минимум 32 символа) |
| `COOKIE_DOMAIN` | No | (empty) | Домен для cookies (например, `localhost`, `.example.com`) |

**Пример:**
```bash
SESSION_SECRET=your-super-long-random-secret-key-at-least-32-chars-long
COOKIE_DOMAIN=localhost  # Для локальной разработки
# или
COOKIE_DOMAIN=.example.com  # Для production (с точкой для поддоменов)
```

**Важные замечания:**
- 🔒 `SESSION_SECRET` **должен быть случайной строкой** минимум 32 символа
- 🔄 Используйте **разные secrets** для dev и production
- 🚫 **Никогда** не коммитьте реальные secrets в Git
- 🎲 Сгенерируйте случайный secret:
  ```bash
  openssl rand -base64 32
  ```

---

## OAuth Configuration

### Base URLs

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OAUTH_BASE_URL` | Yes | - | Публичный URL приложения (где доступно приложение) |
| `OAUTH_FRONTEND_URL` | No | `OAUTH_BASE_URL` | URL для редиректа после OAuth (обычно равен `BASE_URL`) |

**Пример:**
```bash
# Локальная разработка
OAUTH_BASE_URL=http://localhost
OAUTH_FRONTEND_URL=http://localhost

# Production
OAUTH_BASE_URL=https://example.com
OAUTH_FRONTEND_URL=https://example.com
```

**Важно:**
- ✅ URL должен начинаться с `http://` или `https://`
- ✅ В production **обязательно** используйте HTTPS
- ❌ Не добавляйте trailing slash: `https://example.com` (не `https://example.com/`)
- ❌ Не указывайте порт (кроме нестандартных случаев)

### OAuth Providers

#### Google OAuth

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GOOGLE_CLIENT_ID` | No | - | Google OAuth Client ID |
| `GOOGLE_CLIENT_SECRET` | No | - | Google OAuth Client Secret |

**Пример:**
```bash
GOOGLE_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-abcdefghijklmnopqrstuvwx
```

#### GitHub OAuth

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GITHUB_CLIENT_ID` | No | - | GitHub OAuth Client ID |
| `GITHUB_CLIENT_SECRET` | No | - | GitHub OAuth Client Secret |

**Пример:**
```bash
GITHUB_CLIENT_ID=Iv1.a1b2c3d4e5f6g7h8
GITHUB_CLIENT_SECRET=0123456789abcdef0123456789abcdef01234567
```

#### VK ID OAuth

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VK_CLIENT_ID` | No | - | VK Application ID (Client ID) |
| `VK_CLIENT_SECRET` | No | - | VK Secure Key (Client Secret) |

**Пример:**
```bash
VK_CLIENT_ID=54360273
VK_CLIENT_SECRET=ffGI2MwbGCefpcYgxpBy
```

---

## Frontend Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VITE_BACKEND_URL` | No | - | Backend API URL (используется только Vite dev server proxy) |

**Пример:**
```bash
# В Docker Compose (внутренний URL)
VITE_BACKEND_URL=http://backend:8080

# Для локального запуска frontend
VITE_BACKEND_URL=http://localhost:8080
```

**Замечание:** В Docker с Nginx эта переменная не используется, так как все запросы идут через Nginx proxy.

---

## Примеры конфигурации для разных окружений

### 1. Локальная разработка (Docker Compose)

**Файл: `.env` или в `docker-compose.yml`**

```bash
# Core
CONFIG_PATH=config/local.yaml
DATABASE_URL=postgres://postgres:postgres@db:5432/pwp_db?sslmode=disable

# URLs
OAUTH_BASE_URL=http://localhost
OAUTH_FRONTEND_URL=http://localhost
COOKIE_DOMAIN=localhost

# OAuth (опционально, можно оставить в local.yaml)
# GOOGLE_CLIENT_ID=...
# GOOGLE_CLIENT_SECRET=...
# GITHUB_CLIENT_ID=...
# GITHUB_CLIENT_SECRET=...
# VK_CLIENT_ID=...
# VK_CLIENT_SECRET=...
```

**Команды:**
```bash
make dev  # Запуск в Docker Compose
make logs # Просмотр логов
```

---

### 2. Production (HTTPS с реальным доменом)

**Файл: `.env` (НЕ коммитить в Git!)**

```bash
# Core
CONFIG_PATH=config/production.yaml
DATABASE_URL=postgres://produser:securepass@prod-db.internal:5432/production_db?sslmode=require

# Security
SESSION_SECRET=your-super-secure-random-secret-at-least-32-characters-long
COOKIE_DOMAIN=.example.com

# URLs (HTTPS обязательно!)
OAUTH_BASE_URL=https://example.com
OAUTH_FRONTEND_URL=https://example.com

# Google OAuth
GOOGLE_CLIENT_ID=123456789-production-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-production-client-secret

# GitHub OAuth
GITHUB_CLIENT_ID=Iv1.production_client_id
GITHUB_CLIENT_SECRET=production_github_client_secret_40_chars

# VK OAuth
VK_CLIENT_ID=12345678
VK_CLIENT_SECRET=ProductionVKSecureKey
```

**Deployment с Docker:**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

**Deployment с Kubernetes:**
```yaml
# В Secret
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  SESSION_SECRET: "..."
  DATABASE_URL: "..."
  GOOGLE_CLIENT_SECRET: "..."
  # ... остальные secrets

# В ConfigMap (не-секретные данные)
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  OAUTH_BASE_URL: "https://example.com"
  COOKIE_DOMAIN: ".example.com"
  # ...
```

---

### 3. Staging Environment

```bash
# Core
CONFIG_PATH=config/production.yaml
DATABASE_URL=postgres://staginguser:pass@staging-db:5432/staging_db?sslmode=require

# Security (используйте отдельный secret!)
SESSION_SECRET=different-secret-for-staging-environment
COOKIE_DOMAIN=staging.example.com

# URLs
OAUTH_BASE_URL=https://staging.example.com
OAUTH_FRONTEND_URL=https://staging.example.com

# OAuth (можно использовать те же credentials или отдельные)
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
# ...
```

---

## Приоритет конфигурации

Порядок применения настроек (от высшего к низшему приоритету):

1. **Environment variables** (самый высокий приоритет)
2. **YAML config file** (указанный в `CONFIG_PATH`)
3. **Default values** (hardcoded в коде)

**Пример:**

```yaml
# config/local.yaml
oauth:
  base_url: "http://localhost"  # ← Это значение
```

```bash
# Environment variable
export OAUTH_BASE_URL="https://example.com"  # ← Переопределяет YAML
```

**Результат:** Будет использовано `https://example.com` из environment variable.

---

## Проверка конфигурации

### Просмотр используемых значений

```bash
# Запустите приложение и проверьте логи
make dev
make logs

# Ищите строки:
# oauth: initialization complete base_url=http://localhost frontend_url=http://localhost
```

### Тестирование OAuth

```bash
# Backend логи покажут callback URLs для каждого провайдера
make logs | grep "oauth:"

# Пример вывода:
# oauth: Google provider enabled callback_url=http://localhost/auth/google/callback
# oauth: GitHub provider enabled callback_url=http://localhost/auth/github/callback
# oauth: VK ID provider enabled callback_url=http://localhost/auth/vk/callback
# oauth: initialization complete count=3 providers=google, github, vk
```

---

## Troubleshooting

### Ошибка: "config file does not exist"

**Причина:** `CONFIG_PATH` указывает на несуществующий файл.

**Решение:**
```bash
# Проверьте путь
ls -la backend/config/

# Убедитесь что CONFIG_PATH правильный
echo $CONFIG_PATH

# Для Docker Compose:
CONFIG_PATH=config/local.yaml  # Относительный путь от /app в контейнере
```

### Ошибка: "oauth.base_url must start with http:// or https://"

**Причина:** `OAUTH_BASE_URL` не содержит протокол.

**Решение:**
```bash
# Неправильно:
OAUTH_BASE_URL=localhost
OAUTH_BASE_URL=example.com

# Правильно:
OAUTH_BASE_URL=http://localhost
OAUTH_BASE_URL=https://example.com
```

### Ошибка: "cannot read config: ..."

**Причина:** Синтаксическая ошибка в YAML файле или отсутствуют required поля.

**Решение:**
1. Проверьте синтаксис YAML (отступы, двоеточия)
2. Убедитесь что все required переменные установлены:
   - `DATABASE_URL`
   - `SESSION_SECRET`
   - `OAUTH_BASE_URL`

### Environment variables не применяются

**Причина:** Переменные не экспортированы или не переданы в контейнер.

**Решение:**

**Для локального запуска:**
```bash
export OAUTH_BASE_URL=http://localhost
./backend
```

**Для Docker Compose:**
```yaml
# docker-compose.yml
services:
  backend:
    environment:
      - OAUTH_BASE_URL=http://localhost
    # или через .env файл
    env_file:
      - .env
```

**Для Docker run:**
```bash
docker run -e OAUTH_BASE_URL=http://localhost myapp
```

---

## Безопасность Environment Variables

### Best Practices:

1. ✅ **Используйте `.env` файл** для локальной разработки
2. ✅ **Добавьте `.env` в `.gitignore`** (уже добавлен)
3. ✅ **Никогда не коммитьте** secrets в Git
4. ✅ **Используйте разные secrets** для каждого окружения
5. ✅ **Ротируйте secrets** регулярно (особенно если они скомпрометированы)
6. ✅ **В production** используйте секретные хранилища:
   - Docker Secrets
   - Kubernetes Secrets
   - AWS Secrets Manager
   - HashiCorp Vault
   - Azure Key Vault

### Локальная разработка:

```bash
# Создайте .env файл (уже в .gitignore)
cat > .env <<EOF
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pwp_db
SESSION_SECRET=$(openssl rand -base64 32)
OAUTH_BASE_URL=http://localhost
EOF

# Загрузите переменные
source .env

# Или используйте с docker-compose
docker-compose --env-file .env up
```

### Production:

```bash
# НИКОГДА не храните secrets в файлах!
# Используйте CI/CD системы или секретные хранилища

# Пример с Docker Secrets
echo "super-secret-key" | docker secret create session_secret -

# Пример с Kubernetes
kubectl create secret generic app-secrets \
  --from-literal=SESSION_SECRET='your-secret' \
  --from-literal=DATABASE_URL='postgres://...'
```

---

## Дополнительные ресурсы

- [12-Factor App Config](https://12factor.net/config)
- [OWASP Secrets Management](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [Docker Compose Environment Variables](https://docs.docker.com/compose/environment-variables/)
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
