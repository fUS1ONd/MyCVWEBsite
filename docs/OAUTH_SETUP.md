# OAuth Provider Setup Guide

Это руководство объясняет, как зарегистрировать приложение у OAuth провайдеров и получить необходимые credentials.

## Предварительные требования

Для каждого OAuth провайдера вам нужно получить:
- **Client ID** (Application ID)
- **Client Secret** (Secure Key)

## Callback URLs

Все OAuth провайдеры требуют регистрации callback URLs, куда будет перенаправлен пользователь после авторизации.

### Локальная разработка

```
Google:  http://localhost/auth/google/callback
GitHub:  http://localhost/auth/github/callback
VK:      http://localhost/auth/vk/callback
```

### Production

Замените `localhost` на ваш реальный домен:

```
Google:  https://yourdomain.com/auth/google/callback
GitHub:  https://yourdomain.com/auth/github/callback
VK:      https://yourdomain.com/auth/vk/callback
```

---

## 1. VK ID OAuth Setup (VKontakte)

VK ID использует OAuth 2.1 с PKCE (Proof Key for Code Exchange) для дополнительной безопасности.

### Шаги регистрации:

1. **Перейдите в [VK Dev Portal](https://dev.vk.com/)**

2. **Создайте приложение**:
   - Нажмите "Создать приложение"
   - Выберите "Веб-сайт"
   - Укажите название приложения

3. **Настройте OAuth**:
   - Перейдите в **Settings** → **OAuth settings**
   - В поле **Authorized redirect URI** добавьте:
     - Локальная разработка: `http://localhost/auth/vk/callback`
     - Production: `https://yourdomain.com/auth/vk/callback`

4. **Скопируйте credentials**:
   - **Application ID** (Client ID) - в настройках приложения
   - **Secure Key** (Client Secret) - в настройках приложения

5. **Обновите конфигурацию**:
   ```yaml
   # backend/config/local.yaml
   oauth:
     vk:
       client_id: "your-vk-application-id"
       client_secret: "your-vk-secure-key"
       enabled: true
   ```

### Важные замечания для VK:

- ⚠️ VK **НЕ принимает** callback URLs с явным указанием порта (например, `http://localhost:8080`)
- ✅ Используйте Nginx proxy для проксирования на стандартный порт 80/443
- 📝 Callback URL должен **точно совпадать** с зарегистрированным в VK
- 🔒 VK ID поддерживает PKCE для дополнительной безопасности (реализовано автоматически)

---

## 2. Google OAuth Setup

### Шаги регистрации:

1. **Перейдите в [Google Cloud Console](https://console.cloud.google.com/)**

2. **Создайте проект** (или выберите существующий):
   - Нажмите "Select a project" → "New Project"
   - Укажите название проекта
   - Нажмите "Create"

3. **Включите Google+ API**:
   - В меню слева выберите "APIs & Services" → "Library"
   - Найдите "Google+ API"
   - Нажмите "Enable"

4. **Создайте OAuth credentials**:
   - Перейдите в "APIs & Services" → "Credentials"
   - Нажмите "Create Credentials" → "OAuth Client ID"
   - Выберите "Web application"

5. **Настройте OAuth consent screen** (если требуется):
   - User Type: External
   - Заполните обязательные поля (название приложения, email)
   - Добавьте scopes: `email`, `profile`

6. **Настройте Authorized origins**:
   ```
   http://localhost           (для локальной разработки)
   https://yourdomain.com     (для production)
   ```

7. **Настройте Authorized redirect URIs**:
   ```
   http://localhost/auth/google/callback
   https://yourdomain.com/auth/google/callback
   ```

8. **Скопируйте credentials**:
   - Client ID
   - Client Secret

9. **Обновите конфигурацию**:
   ```yaml
   # backend/config/local.yaml
   oauth:
     google:
       client_id: "your-client-id.apps.googleusercontent.com"
       client_secret: "your-client-secret"
       enabled: true
   ```

---

## 3. GitHub OAuth Setup

### Шаги регистрации:

1. **Перейдите в [GitHub Settings](https://github.com/settings/developers)**
   - Developer Settings → OAuth Apps

2. **Создайте новое OAuth приложение**:
   - Нажмите "New OAuth App"

3. **Заполните форму**:
   - **Application name**: Название вашего приложения
   - **Homepage URL**:
     - Локально: `http://localhost`
     - Production: `https://yourdomain.com`
   - **Application description**: (опционально) Описание приложения
   - **Authorization callback URL**:
     - Локально: `http://localhost/auth/github/callback`
     - Production: `https://yourdomain.com/auth/github/callback`

4. **Зарегистрируйте приложение**:
   - Нажмите "Register application"

5. **Скопируйте Client ID**:
   - Отображается сразу после создания

6. **Сгенерируйте Client Secret**:
   - Нажмите "Generate a new client secret"
   - **ВАЖНО**: Скопируйте secret сразу, он показывается только один раз!

7. **Обновите конфигурацию**:
   ```yaml
   # backend/config/local.yaml
   oauth:
     github:
       client_id: "your-github-client-id"
       client_secret: "your-github-client-secret"
       enabled: true
   ```

---

## Настройка для разных окружений

### Локальная разработка (Docker Compose)

Отредактируйте `backend/config/local.yaml`:

```yaml
oauth:
  base_url: "http://localhost"
  frontend_url: "http://localhost"

  google:
    client_id: "your-google-client-id"
    client_secret: "your-google-client-secret"
    enabled: true

  github:
    client_id: "your-github-client-id"
    client_secret: "your-github-client-secret"
    enabled: true

  vk:
    client_id: "your-vk-client-id"
    client_secret: "your-vk-client-secret"
    enabled: true
```

### Production (через Environment Variables)

Используйте environment variables для безопасного хранения credentials:

```bash
export OAUTH_BASE_URL="https://yourdomain.com"
export OAUTH_FRONTEND_URL="https://yourdomain.com"

export GOOGLE_CLIENT_ID="xxx.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="xxx"

export GITHUB_CLIENT_ID="xxx"
export GITHUB_CLIENT_SECRET="xxx"

export VK_CLIENT_ID="xxx"
export VK_CLIENT_SECRET="xxx"
```

Или используйте `.env` файл, docker secrets, Kubernetes secrets, etc.

---

## Тестирование OAuth провайдеров

После настройки всех провайдеров:

1. **Перезапустите приложение**:
   ```bash
   make reset
   make dev
   ```

2. **Проверьте логи**:
   ```bash
   make logs
   ```
   Вы должны увидеть:
   ```
   oauth: Google provider enabled callback_url=http://localhost/auth/google/callback
   oauth: GitHub provider enabled callback_url=http://localhost/auth/github/callback
   oauth: VK ID provider enabled callback_url=http://localhost/auth/vk/callback
   oauth: initialization complete count=3 providers=google, github, vk
   ```

3. **Тестируйте каждый провайдер**:
   - Откройте `http://localhost` в браузере
   - Нажмите "Login"
   - Попробуйте каждую кнопку OAuth (Google, GitHub, VK)
   - Проверьте, что после авторизации вы перенаправляетесь на `/blog`

---

## Troubleshooting (Решение проблем)

### Ошибка: "Redirect URI mismatch"

**Причина**: Callback URL в конфигурации не совпадает с зарегистрированным в провайдере.

**Решение**:
1. Проверьте логи backend для точного callback URL
2. Убедитесь, что этот URL зарегистрирован в OAuth провайдере
3. Callback URL должен совпадать **полностью** (включая протокол, домен, путь)

### Ошибка: "OAuth callback failed"

**Причина**: Неверные Client ID или Client Secret.

**Решение**:
1. Перепроверьте Client ID и Secret в конфигурации
2. Убедитесь, что OAuth провайдер активен (`enabled: true`)
3. Проверьте логи backend для детальной ошибки

### VK OAuth не работает с портом :8080

**Причина**: VK не принимает callback URLs с нестандартными портами.

**Решение**:
1. Используйте Nginx proxy (уже настроен в docker-compose)
2. Callback URL должен быть `http://localhost/auth/vk/callback` (без :8080)
3. В `base_url` также не должно быть :8080

### Cookie не отправляется после OAuth

**Причина**: Cookie domain не соответствует домену в браузере.

**Решение**:
1. Проверьте `cookie_domain` в config
2. Для localhost используйте `cookie_domain: "localhost"`
3. Для production используйте ваш домен (например, `"example.com"`)

---

## Безопасность

### Best Practices:

1. ✅ **Никогда не коммитьте** реальные Client Secret в Git
2. ✅ **Используйте разные credentials** для dev и production
3. ✅ **Ротируйте secrets** регулярно
4. ✅ **В production** всегда используйте HTTPS (`cookie_secure: true`)
5. ✅ **Используйте секретные хранилища** (Vault, AWS Secrets Manager) в production

### Для локальной разработки:

- Можно использовать `local.yaml` с реальными credentials
- Добавьте `local.yaml` в `.gitignore` (уже добавлен)
- Или используйте `.env` файл (тоже в .gitignore)

### Для production:

- **Обязательно** используйте environment variables
- Не храните secrets в config файлах
- Используйте `production.yaml.example` как шаблон

---

## Дополнительные ресурсы

- [VK ID Documentation](https://id.vk.ru/about/business/go/docs)
- [Google OAuth 2.0 Guide](https://developers.google.com/identity/protocols/oauth2)
- [GitHub OAuth Guide](https://docs.github.com/en/developers/apps/building-oauth-apps)
- [OWASP OAuth Security](https://cheatsheetseries.owasp.org/cheatsheets/OAuth2_Cheat_Sheet.html)
