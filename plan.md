# Детальный План Разработки: Персональный веб-сайт (CV + AI Blog)

**Проект:** Персональный веб-сайт с CV-страницей и AI Blog
**Стек:** Go 1.22+, Chi, pgx/v5, React 18+, TypeScript, PostgreSQL 15
**Дата создания плана:** 2025-11-25
**Версия плана:** 1.0

---

## Выбранный стек технологий

### Backend
- **Язык**: Go 1.22+
- **HTTP Router**: Chi v5 (уже используется)
- **БД**: PostgreSQL 15
- **ORM/Driver**: pgx/v5 (уже используется)
- **Миграции**: golang-migrate
- **OAuth**: goth (VK, Google, GitHub)
- **Логирование**: slog (уже используется)
- **Валидация**: go-playground/validator
- **Тестирование**: testify, testcontainers

### Frontend
- **Framework**: React 18+ с TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **State Management**: React Query + Context API
- **Rich Text Editor**: TipTap или Draft.js
- **Тестирование**: Vitest, React Testing Library
- **E2E**: Playwright

### Infrastructure
- **Контейнеризация**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **SSL**: Let's Encrypt / Certbot
- **File Storage**: локальное хранилище на сервере
- **Notifications**: SMTP + фоновые горутины
- **Deploy**: VPS (DigitalOcean/Hetzner)

### Аналитика и мониторинг
- Яндекс.Метрика
- Google Analytics 4
- Sentry (опционально для error tracking)

---

## Текущий статус проекта

**Уже реализовано:**
- ✅ Базовая структура Go проекта
- ✅ Chi router
- ✅ PostgreSQL подключение (pgx/v5)
- ✅ Docker Compose с PostgreSQL
- ✅ Базовая миграция users
- ✅ Clean Architecture (domain, repository, service, transport)
- ✅ Логирование через slog

**Требуется реализовать:**
- OAuth интеграция (goth)
- Frontend полностью (React + Vite + TypeScript + Tailwind)
- Миграции для posts, comments, media, notifications
- Rich text editor для админ-панели
- File upload система
- Email уведомления
- SEO и Analytics интеграция
- Тестирование (70%+ backend, 60%+ frontend)
- Production деплой с CI/CD

---

## Этап 1: Инфраструктура и базовая настройка

**Цель этапа:** Настроить полную инфраструктуру проекта, системы сборки, линтеры, CI/CD базу

**Список задач:**

- [x] Настроить Makefile с командами (build, test, lint, migrate, docker-up, docker-down)
- [x] Добавить .gitignore для Go и Node.js проектов
- [x] Настроить Go modules и добавить необходимые зависимости (chi, pgx, goth, validator, testify)
- [x] Создать структуру директорий для frontend (монорепозиторий: backend/ и frontend/)
- [x] Инициализировать React проект с Vite и TypeScript
- [x] Настроить Tailwind CSS в frontend проекте
- [x] Добавить pre-commit хуки (go fmt, go vet, eslint, prettier)
- [x] Настроить ESLint и Prettier для TypeScript/React
- [x] Создать README.md с инструкциями по запуску и разработке
- [x] Настроить golangci-lint с конфигурацией
- [x] Настроить CI/CD (GitHub Actions для backend и frontend)
- [x] Обновить docker-compose.yml для разработки (hot reload для Go и React)
- [x] Добавить healthcheck endpoints (/health, /ready)

**Критерии завершения:**
- Проект собирается одной командой (make build)
- Docker compose поднимает всю инфраструктуру
- Линтеры работают и проходят без ошибок
- Hot reload работает для backend и frontend

---

## Этап 2: База данных и модели

**Цель этапа:** Спроектировать и создать все необходимые таблицы БД, настроить миграции

**Список задач:**

- [x] Создать миграцию для таблицы oauth_providers (id, user_id, provider, provider_user_id, access_token, refresh_token, expires_at)
- [x] Создать миграцию для таблицы sessions (id, user_id, token, expires_at, created_at)
- [x] Создать миграцию для таблицы posts (id, title, slug, content, preview, author_id, published, published_at, created_at, updated_at)
- [x] Создать миграцию для таблицы comments (id, post_id, user_id, content, parent_id, created_at, updated_at, deleted_at)
- [x] Создать миграцию для таблицы media_files (id, filename, path, mime_type, size, uploader_id, created_at)
- [x] Создать миграцию для таблицы post_media (post_id, media_id, sort_order)
- [x] Создать миграцию для таблицы notification_settings (user_id, email_enabled, push_enabled, new_posts_enabled)
- [x] Создать миграцию для таблицы notifications (id, user_id, type, title, message, read, created_at)
- [x] Добавить индексы для оптимизации (posts.slug, posts.published_at, comments.post_id, sessions.token)
- [x] Создать Domain модели для всех сущностей в internal/domain/
- [x] Добавить валидацию на уровне domain models (go-playground/validator)
- [x] Создать seed-скрипт для тестовых данных (admin user, sample posts)

**Критерии завершения:**
- Все миграции применяются без ошибок
- Rollback миграций работает корректно
- Domain модели покрывают все требования ТЗ
- Seed данные создают полноценное окружение для разработки

---

## Этап 3: Аутентификация и авторизация

**Цель этапа:** Реализовать OAuth аутентификацию через VK, Google, GitHub и систему сессий

**Список задач:**

- [x] Установить и настроить библиотеку goth для OAuth
- [x] Создать структуру config для OAuth провайдеров (client_id, client_secret, callback_url)
- [x] Реализовать AuthRepository (CreateUser, GetUserByEmail, LinkOAuthProvider, GetUserByProviderID)
- [x] Реализовать SessionRepository (CreateSession, GetSession, DeleteSession, CleanupExpiredSessions)
- [x] Создать AuthService с методами (LoginWithOAuth, Logout, ValidateSession, RefreshSession)
- [x] Реализовать OAuth endpoints (GET /auth/{provider}, GET /auth/{provider}/callback)
- [x] Создать middleware для проверки аутентификации (AuthRequired)
- [x] Создать middleware для проверки роли администратора (AdminRequired)
- [x] Реализовать механизм хранения JWT токенов или session ID в cookies (httpOnly, secure)
- [ ] Добавить CSRF protection middleware
- [x] Создать endpoint GET /auth/me для получения информации о текущем пользователе
- [x] Создать endpoint POST /auth/logout для выхода
- [x] Добавить автоматическую очистку просроченных сессий (background goroutine)
- [ ] Написать unit тесты для AuthService
- [ ] Написать integration тесты для auth endpoints

**Критерии завершения:**
- Авторизация через VK, Google и GitHub работает
- Сессии корректно создаются и валидируются
- Middleware защищает закрытые endpoints
- CSRF protection работает
- Тесты покрывают > 70% кода auth модуля

---

## Этап 4: Backend API

**Цель этапа:** Создать REST API для всех функций приложения

**Список задач:**

### 4.1 Profile API (публичный)
- [x] Реализовать ProfileRepository с методом GetProfile() из БД (вместо config)
- [x] Создать миграцию для таблицы profile_info (id, name, description, photo_url, activity, updated_at)
- [x] Реализовать ProfileService (GetProfile, UpdateProfile - только для админа)
- [x] Создать endpoint GET /api/v1/profile (публичный)
- [x] Создать endpoint PUT /api/v1/admin/profile (требует AdminRequired)

### 4.2 Posts API
- [x] Реализовать PostRepository (Create, Update, Delete, GetByID, GetBySlug, List with pagination)
- [x] Реализовать PostService с бизнес-логикой (валидация, slug generation, permission checks)
- [x] Создать endpoint GET /api/v1/posts (список, пагинация, фильтры по published)
- [x] Создать endpoint GET /api/v1/posts/{slug} (детали поста)
- [x] Создать endpoint POST /api/v1/admin/posts (создание, AdminRequired)
- [x] Создать endpoint PUT /api/v1/admin/posts/{id} (редактирование, AdminRequired)
- [x] Создать endpoint DELETE /api/v1/admin/posts/{id} (удаление, AdminRequired)
- [x] Добавить поддержку draft/published статусов
- [x] Реализовать автогенерацию slug из title (транслитерация)

### 4.3 Comments API
- [ ] Реализовать CommentRepository (Create, Update, Delete, GetByPostID, GetByID)
- [ ] Реализовать CommentService (валидация, permission checks, nested comments support)
- [ ] Создать endpoint GET /api/v1/posts/{slug}/comments (список комментариев с вложенностью)
- [ ] Создать endpoint POST /api/v1/posts/{slug}/comments (создание, AuthRequired)
- [ ] Создать endpoint PUT /api/v1/comments/{id} (редактирование своего комментария)
- [ ] Создать endpoint DELETE /api/v1/comments/{id} (удаление своего или AdminRequired)
- [ ] Реализовать soft delete для комментариев (deleted_at)

### 4.4 Общие задачи API
- [ ] Добавить middleware для логирования всех запросов
- [ ] Добавить middleware для CORS (настраиваемые origins)
- [ ] Реализовать стандартизированные JSON responses (success/error format)
- [ ] Добавить валидацию request body с понятными error messages
- [ ] Реализовать rate limiting middleware (по IP и по user_id)
- [ ] Добавить middleware для panic recovery
- [ ] Создать OpenAPI/Swagger документацию для API
- [ ] Написать integration тесты для всех endpoints

**Критерии завершения:**
- Все CRUD операции работают корректно
- Пагинация и фильтры работают
- Авторизация и валидация работает на всех endpoints
- API документирован
- Integration тесты покрывают все endpoints

---

## Этап 5: Frontend основа

**Цель этапа:** Настроить базовую архитектуру React приложения, роутинг, API клиент

**Список задач:**

- [ ] Настроить React Router v6 с route конфигурацией
- [ ] Создать структуру папок (components/, pages/, hooks/, services/, types/, utils/)
- [ ] Создать API client с axios (base URL, interceptors, error handling)
- [ ] Реализовать auth context (AuthProvider, useAuth hook)
- [ ] Создать protected route компонент (требует авторизации)
- [ ] Создать admin route компонент (требует роли admin)
- [ ] Настроить TypeScript types для всех API моделей (User, Post, Comment, Profile)
- [ ] Создать базовые UI компоненты (Button, Input, Card, Modal, Spinner)
- [ ] Настроить общий Layout компонент с Header и Footer
- [ ] Создать ErrorBoundary для обработки ошибок React
- [ ] Настроить react-query для кэширования и управления server state
- [ ] Создать toast notification system для user feedback
- [ ] Добавить Loading states и Error states для всех асинхронных операций

**Критерии завершения:**
- Роутинг работает корректно
- API клиент успешно делает запросы к backend
- Auth context корректно управляет состоянием авторизации
- Базовые UI компоненты переиспользуемы и типизированы
- React Query оптимизирует запросы к API

---

## Этап 6: Публичная CV страница

**Цель этапа:** Реализовать публичную CV-страницу со всей необходимой информацией

**Список задач:**

- [ ] Создать HomePage компонент (/), доступный без авторизации
- [ ] Реализовать Hero section с фотографией и кратким описанием
- [ ] Создать About section (о себе, навыки, описание)
- [ ] Создать Activity section (текущая деятельность)
- [ ] Создать Contacts section (email, GitHub, LinkedIn, VK)
- [ ] Добавить адаптивную верстку для мобильных устройств
- [ ] Реализовать smooth scroll между секциями
- [ ] Добавить animations/transitions для улучшения UX
- [ ] Оптимизировать загрузку фотографии (webp format, lazy load)
- [ ] Добавить meta tags для SEO (title, description, keywords)
- [ ] Добавить Open Graph meta tags для социальных сетей
- [ ] Реализовать темную/светлую тему (опционально, но желательно)
- [ ] Добавить кнопку "Войти" в Header для перехода к OAuth

**Критерии завершения:**
- CV страница выглядит профессионально
- Все секции заполнены актуальной информацией
- Адаптивная верстка работает на всех устройствах
- Страница загружается < 2 секунд
- SEO meta tags правильно настроены

---

## Этап 7: Система постов и комментариев

**Цель этапа:** Реализовать функционал просмотра постов и комментирования для авторизованных пользователей

**Список задач:**

### 7.1 Страница списка постов
- [ ] Создать PostsListPage компонент (/blog), требует авторизации
- [ ] Реализовать PostCard компонент (превью поста)
- [ ] Добавить пагинацию для списка постов
- [ ] Реализовать фильтры (по дате, поиск по заголовку)
- [ ] Добавить skeleton loading для постов
- [ ] Оптимизировать изображения в превью (thumbnails)

### 7.2 Страница поста
- [ ] Создать PostDetailPage компонент (/blog/{slug})
- [ ] Реализовать рендеринг markdown контента (react-markdown + syntax highlighting)
- [ ] Добавить поддержку встроенных изображений в контенте
- [ ] Добавить поддержку встроенных видео (YouTube, Vimeo embeds)
- [ ] Реализовать share кнопки (Twitter, Facebook, LinkedIn, VK)
- [ ] Добавить кнопку "Назад к списку постов"
- [ ] Реализовать навигацию к предыдущему/следующему посту

### 7.3 Комментарии
- [ ] Создать CommentsSection компонент
- [ ] Реализовать Comment компонент (аватар, имя, дата, текст)
- [ ] Добавить форму для создания комментария (textarea + submit)
- [ ] Реализовать вложенные комментарии (reply to comment)
- [ ] Добавить возможность редактирования своего комментария
- [ ] Добавить возможность удаления своего комментария (с подтверждением)
- [ ] Реализовать optimistic updates для комментариев
- [ ] Добавить валидацию комментариев на frontend (max length, required)
- [ ] Реализовать auto-scroll к новому комментарию после создания

**Критерии завершения:**
- Авторизованные пользователи видят все посты
- Markdown контент корректно отображается
- Изображения и видео встраиваются в посты
- Комментарии работают с вложенностью
- Пользователи могут редактировать/удалять свои комментарии
- Пагинация работает корректно

---

## Этап 8: Административная панель

**Цель этапа:** Создать полноценную админ-панель для управления контентом

**Список задач:**

### 8.1 Структура админки
- [ ] Создать AdminLayout компонент с боковым меню
- [ ] Создать Dashboard страницу (/admin) со статистикой (кол-во постов, комментариев, users)
- [ ] Добавить admin route защиту для всех админ страниц

### 8.2 Управление постами
- [ ] Создать AdminPostsListPage (/admin/posts) с таблицей всех постов
- [ ] Добавить фильтры (published/draft, дата, поиск)
- [ ] Создать CreatePostPage (/admin/posts/new)
- [ ] Создать EditPostPage (/admin/posts/{id}/edit)
- [ ] Интегрировать Rich Text Editor (TipTap или Draft.js с markdown поддержкой)
- [ ] Добавить preview режим для поста перед публикацией
- [ ] Реализовать auto-save для черновиков (debounced)
- [ ] Добавить возможность загрузки изображений прямо в редакторе
- [ ] Реализовать управление slug (auto-generate или manual)
- [ ] Добавить поля: title, slug, content, preview, published status
- [ ] Создать DeletePost функционал с подтверждением

### 8.3 Модерация комментариев
- [ ] Создать AdminCommentsPage (/admin/comments)
- [ ] Добавить фильтры (по посту, дате, пользователю)
- [ ] Реализовать массовое удаление комментариев
- [ ] Добавить возможность редактирования комментариев администратором
- [ ] Создать модалку для просмотра полного контекста комментария

### 8.4 Управление профилем
- [ ] Создать AdminProfileEditPage (/admin/profile)
- [ ] Добавить форму редактирования CV информации
- [ ] Реализовать загрузку/замену фотографии профиля
- [ ] Добавить preview изменений перед сохранением

**Критерии завершения:**
- Администратор может создавать/редактировать/удалять посты
- Rich Text Editor работает корректно с markdown
- Изображения загружаются и встраиваются в посты
- Модерация комментариев работает
- CV информация редактируется через админку
- Auto-save защищает от потери данных

---

## Этап 9: Система уведомлений

**Цель этапа:** Реализовать отправку email и push уведомлений о новых постах

**Список задач:**

### 9.1 Email уведомления
- [ ] Добавить SMTP конфигурацию в config (host, port, username, password, from)
- [ ] Создать EmailService с методом SendEmail(to, subject, body, htmlBody)
- [ ] Создать email templates (новый пост, welcome email)
- [ ] Реализовать NotificationRepository (CreateNotification, MarkAsRead, GetUserNotifications)
- [ ] Создать NotificationService с методом NotifyNewPost(postID)
- [ ] Интегрировать отправку уведомлений при публикации поста
- [ ] Реализовать фоновую горутину для асинхронной отправки email
- [ ] Добавить retry механизм при ошибках отправки

### 9.2 Push уведомления
- [ ] Настроить service worker для push notifications
- [ ] Реализовать запрос разрешения на push в браузере
- [ ] Создать endpoint для регистрации push subscription
- [ ] Реализовать хранение push subscriptions в БД
- [ ] Интегрировать отправку push уведомлений при публикации поста
- [ ] Добавить обработку click на push notification (redirect на пост)

### 9.3 Настройки уведомлений
- [ ] Создать NotificationSettingsPage (/settings/notifications)
- [ ] Добавить переключатели (email, push, новые посты)
- [ ] Реализовать сохранение настроек пользователя
- [ ] Учитывать настройки при отправке уведомлений
- [ ] Создать страницу истории уведомлений (/notifications)

**Критерии завершения:**
- Email уведомления отправляются корректно
- Push уведомления работают в браузере
- Пользователи могут управлять своими настройками уведомлений
- Уведомления отправляются асинхронно без блокировки API
- Retry механизм обрабатывает временные ошибки

---

## Этап 10: File Upload и медиа

**Цель этапа:** Реализовать загрузку и управление файлами (изображения, видео)

**Список задач:**

- [ ] Создать директорию для хранения uploads на сервере (./storage/uploads)
- [ ] Реализовать MediaRepository (Create, Delete, GetByID, GetByUploader)
- [ ] Создать MediaService с валидацией (mime types, file size, dimensions)
- [ ] Создать endpoint POST /api/v1/admin/media (загрузка файлов, multipart/form-data)
- [ ] Создать endpoint DELETE /api/v1/admin/media/{id} (удаление файлов)
- [ ] Создать endpoint GET /api/v1/media/{filename} (раздача файлов со статикой)
- [ ] Добавить генерацию thumbnails для изображений (resize, webp conversion)
- [ ] Реализовать ограничения по размеру файла (max 10MB для изображений, 50MB для видео)
- [ ] Добавить ограничения по типам файлов (jpg, png, webp, gif для изображений; mp4, webm для видео)
- [ ] Создать FileUpload React компонент с drag-and-drop
- [ ] Добавить progress bar для загрузки файлов
- [ ] Реализовать preview загруженных изображений
- [ ] Создать MediaLibrary компонент для выбора существующих файлов
- [ ] Интегрировать загрузку файлов в Rich Text Editor
- [ ] Добавить автоматическую очистку неиспользуемых файлов (background job)

**Критерии завершения:**
- Загрузка изображений и видео работает
- Thumbnails генерируются автоматически
- Drag-and-drop работает в админке
- Файлы корректно раздаются через API
- Валидация защищает от загрузки недопустимых файлов
- MediaLibrary позволяет переиспользовать загруженные файлы

---

## Этап 11: SEO и аналитика

**Цель этапа:** Оптимизировать сайт для поисковых систем и интегрировать аналитику

**Список задач:**

### 11.1 SEO оптимизация
- [ ] Настроить server-side rendering (SSR) или static site generation для публичных страниц
- [ ] Добавить react-helmet-async для управления meta tags
- [ ] Реализовать динамические meta tags для каждой страницы поста
- [ ] Добавить Open Graph tags (og:title, og:description, og:image, og:url)
- [ ] Добавить Twitter Card tags
- [ ] Создать sitemap.xml endpoint (/sitemap.xml)
- [ ] Создать robots.txt endpoint (/robots.txt)
- [ ] Добавить canonical URLs для всех страниц
- [ ] Реализовать structured data (JSON-LD) для постов (Article schema)
- [ ] Оптимизировать заголовки (H1, H2, H3) в контенте
- [ ] Добавить alt атрибуты для всех изображений
- [ ] Реализовать 404 страницу с навигацией

### 11.2 Аналитика
- [ ] Интегрировать Яндекс.Метрику (добавить счетчик в HTML)
- [ ] Интегрировать Google Analytics 4 (gtag.js)
- [ ] Настроить отслеживание событий (просмотры постов, комментарии, авторизация)
- [ ] Создать AdminAnalyticsPage (/admin/analytics) с виджетами метрик
- [ ] Добавить privacy policy страницу с информацией о cookies

**Критерии завершения:**
- Sitemap.xml генерируется и доступен
- Robots.txt настроен корректно
- Meta tags заполнены для всех страниц
- Яндекс.Метрика и Google Analytics отслеживают события
- Structured data валидна (проверить через Google Rich Results Test)

---

## Этап 12: Тестирование

**Цель этапа:** Достичь требуемого покрытия тестами (Backend 70%+, Frontend 60%+)

**Список задач:**

### 12.1 Backend тестирование
- [ ] Написать unit тесты для всех Service слоев
- [ ] Написать unit тесты для domain models валидации
- [ ] Написать unit тесты для utility функций
- [ ] Написать integration тесты для всех Repository методов (с testcontainers)
- [ ] Написать integration тесты для всех HTTP endpoints
- [ ] Написать тесты для middleware (auth, CSRF, rate limiting)
- [ ] Добавить table-driven tests где применимо
- [ ] Настроить coverage report (go test -cover)
- [ ] Достичь минимум 70% покрытия кода

### 12.2 Frontend тестирование
- [ ] Написать unit тесты для utility функций и hooks
- [ ] Написать component тесты с React Testing Library для UI компонентов
- [ ] Написать integration тесты для страниц (HomePage, PostsListPage, PostDetailPage)
- [ ] Написать тесты для auth flow (login, logout, protected routes)
- [ ] Написать тесты для форм (валидация, отправка)
- [ ] Настроить MSW (Mock Service Worker) для мокирования API
- [ ] Настроить coverage report (vitest --coverage)
- [ ] Достичь минимум 60% покрытия компонентов

### 12.3 E2E тестирование
- [ ] Настроить Playwright или Cypress
- [ ] Написать E2E тест для OAuth авторизации
- [ ] Написать E2E тест для создания поста админом
- [ ] Написать E2E тест для комментирования поста пользователем
- [ ] Написать E2E тест для редактирования профиля админом
- [ ] Написать E2E тест для подписки на уведомления
- [ ] Настроить запуск E2E тестов в CI

**Критерии завершения:**
- Backend покрытие > 70%
- Frontend покрытие > 60%
- Все критические сценарии покрыты E2E тестами
- Тесты проходят стабильно
- CI автоматически запускает все тесты

---

## Этап 13: Деплой и CI/CD

**Цель этапа:** Настроить автоматический деплой на production сервер

**Список задач:**

### 13.1 Production конфигурация
- [ ] Создать production Dockerfile (multi-stage build для Go и React)
- [ ] Создать production docker-compose.yml
- [ ] Настроить environment variables для production
- [ ] Создать systemd unit file для автозапуска (опционально)
- [ ] Настроить логирование в файлы (rotation)

### 13.2 Nginx и SSL
- [ ] Установить и настроить Nginx как reverse proxy
- [ ] Настроить Nginx конфигурацию для раздачи статики
- [ ] Установить и настроить Certbot для Let's Encrypt
- [ ] Получить SSL сертификат
- [ ] Настроить автоматическое обновление сертификата
- [ ] Настроить HTTP -> HTTPS редирект
- [ ] Добавить security headers (HSTS, X-Frame-Options, CSP)
- [ ] Настроить gzip/brotli compression

### 13.3 CI/CD
- [ ] Создать GitHub Actions workflow для тестов (на каждый push)
- [ ] Создать GitHub Actions workflow для линтеров
- [ ] Создать GitHub Actions workflow для build и деплоя (на push в main)
- [ ] Настроить secrets в GitHub (SSH keys, production env vars)
- [ ] Реализовать zero-downtime deployment (blue-green или rolling update)
- [ ] Настроить health checks перед переключением трафика
- [ ] Добавить rollback механизм при ошибках деплоя

### 13.4 Мониторинг
- [ ] Настроить логирование ошибок (Sentry или аналог)
- [ ] Настроить алерты на критические ошибки
- [ ] Добавить мониторинг ресурсов сервера (CPU, RAM, Disk)
- [ ] Настроить автоматический backup PostgreSQL базы данных

### 13.5 Финальная проверка
- [ ] Проверить все функции на production
- [ ] Проверить SSL сертификат (SSL Labs A+ rating)
- [ ] Проверить скорость загрузки (< 3 секунд на всех страницах)
- [ ] Проверить адаптивную верстку на всех устройствах
- [ ] Проверить OAuth авторизацию через все провайдеры
- [ ] Проверить отправку email уведомлений
- [ ] Проверить push уведомления
- [ ] Запустить security audit (проверить OWASP Top 10)

**Критерии завершения:**
- Сайт доступен по HTTPS
- CI/CD автоматически деплоит изменения
- SSL рейтинг A+ на SSL Labs
- Все метрики успеха из ТЗ достигнуты
- Мониторинг и алерты настроены
- Backup базы данных работает

---

## Дополнительные задачи (опционально, после основных этапов)

- [ ] Реализовать полнотекстовый поиск по постам (PostgreSQL full-text search)
- [ ] Добавить систему тегов для постов
- [ ] Реализовать RSS feed для постов
- [ ] Добавить поддержку internationalization (i18n) для русского и английского
- [ ] Реализовать систему лайков для постов и комментариев
- [ ] Добавить возможность прикрепления files к комментариям
- [ ] Реализовать email subscription для уведомлений без авторизации
- [ ] Добавить страницу "About this project" с использованными технологиями
- [ ] Создать API rate limiting dashboard для админа
- [ ] Добавить возможность экспорта постов в Markdown файлы

---

## Метрики успеха проекта

По завершении всех этапов проект должен соответствовать следующим критериям:

- [ ] Сайт загружается за < 3 секунд
- [ ] Адаптивная верстка работает на мобильных, планшетах и десктопах
- [ ] SSL сертификат настроен корректно (A+ рейтинг SSL Labs)
- [ ] OAuth авторизация через VK, Google и GitHub работает
- [ ] Backend покрытие тестами > 70%
- [ ] Frontend покрытие тестами > 60%
- [ ] Нет критических уязвимостей безопасности
- [ ] Яндекс.Метрика и Google Analytics корректно отслеживают события
- [ ] Email и push уведомления работают
- [ ] Админ может управлять всем контентом через удобную панель
- [ ] Пользователи могут комментировать посты
- [ ] Все требования из ТЗ выполнены

---

## Критические файлы для реализации

Наиболее критичные файлы для реализации проекта:

1. **cmd/app/main.go** - Точка входа приложения, инициализация всех компонентов
2. **internal/domain/user.go** - Domain модель пользователя, основа для OAuth
3. **internal/repository/postgres.go** - Подключение к БД и базовая структура репозиториев
4. **config/config.go** - Конфигурация всего приложения (OAuth, SMTP, пути к файлам)
5. **docker-compose.yml** - Инфраструктура для разработки и production

---

**Конец плана разработки**

Этот план полностью покрывает все требования из ТЗ (CLAUDE.md) и структурирован для последовательной разработки от инфраструктуры к финальному деплою. Каждый этап можно выполнять независимо, отмечая прогресс через чекбоксы.
