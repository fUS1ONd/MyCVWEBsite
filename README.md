# Personal Web Platform

Персональный веб-сайт (CV + AI Blog).

## 🚀 Quick Start

**Prerequisites:** Docker, Make.

1. **Initialize Environment** (First time only)

   ```bash
   make init
   ```

2. **Start Development**
   ```bash
   make dev
   ```
   - **Backend API**: http://localhost:8080
   - **Frontend**: http://localhost:5173

## 🛠 Commands

| Command      | Description                                 |
| ------------ | ------------------------------------------- |
| `make init`  | Setup local environment (git hooks, config) |
| `make dev`   | Start everything (Docker + Hot Reload)      |
| `make logs`  | View server logs                            |
| `make stop`  | Stop containers                             |
| `make reset` | Stop and wipe database/volumes              |
| `make check` | Run linters and formatters                  |

docker compose exec -T db psql -U postgres -d pwp_db -c "UPDATE users SET role = 'admin' WHERE email = 'koskriv2006@gmail.com'"

docker compose run --rm migrator

/home/krivonosov/.claude/plans/splendid-singing-blanket.md

Проверь последние 10 коммитов, чтобы проверить, что уже было сделано, согласно плану - /home/krivonosov/.claude/plans/splendid-singing-blanket.md переиспользуй код, который был добавлен в этих коммитах. Референс на примерный интерфейс просмотра постов(хаброподобный)- curriculum_vitae/image.png. Продолжи выполнение плана(вроде четвертый этап). не забывай после выполнения большого мажорного этапа делать make check и git commit, а ещё убеждайся что не делаешь лажу: анализируй код, чтобы не делать одни и те же методы но с разным названием.
