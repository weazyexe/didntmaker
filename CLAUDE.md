# CLAUDE.md

Guidance for AI agents (and humans) working in this repo. `AGENTS.md` and `GEMINI.md` are symlinks to this file.

## What this is

**Несправлятор** — a Telegram bot for group chats. It tracks a chat's *reputation of failing* via points.

The core unit is **очки несправления** ("non-coping points"). The mechanic: reply to someone's message
with `+N` or `-N` to grant or deduct points. Each member has a daily allowance to spend. Extras:

- **/bet** — once-a-day dice roll: 4–6 wins +1000 allowance, 1–3 loses 1000 score.
- **Mass grant** — reply `+N`/`-N` to a *bot* message and it spreads to everyone in the chat.
- **Stats** — leaderboard (`/stats`), per-user stats (`/info`), bet stats (`/bet_stats`).
- **Discord** — bind a guild and get voice-channel join/leave notifications.

The whole product voice is an aggressive lowercase Russian meme tone. See the **i18n** section and the
`i18n-tone` skill before writing any user-facing string.

## Tech stack

- **Go 1.25**
- **telebot.v3** (`gopkg.in/telebot.v3`) — Telegram Bot API
- **discordgo** — Discord API
- **sqlc** — type-safe DB access, **generated** from SQL (not hand-written, not committed)
- **goose** — embedded migrations, applied on startup
- **modernc.org/sqlite** — pure-Go SQLite driver (no cgo)
- **Task** (`Taskfile.yml`) — task runner
- **Docker** — deployment

## Architecture & layering

Strict one-directional dependency flow. Do not skip layers.

```
cmd/{bot,migrate}          entry points, wiring, graceful shutdown
        │
internal/bot               telebot setup + middleware chain
        │
internal/handlers          one file per command; thin, parses input, calls services
        │
internal/service           business logic (transfers, bets, stats, time rules)
        │
internal/repository        maps domain types ↔ generated sqlc code
        │
internal/database/gen      sqlc OUTPUT — generated, gitignored, never edit by hand
```

Cross-cutting packages:

- `internal/domain` — plain types, enums, errors. No external deps. Crosses all layers.
- `internal/i18n` — every user-facing string. One `Messages` struct, one RU impl.
- `internal/config` — env config via cleanenv.
- `internal/database` — opens SQLite, runs migrations.
- `internal/integration` — black-box tests over the real stack.

## The ledger (source of truth)

Points and limits are **not** stored as running totals. The append-only `postings` table
(`db/schema.sql`) is the source of truth; balances are `SUM(amount)`.

Each posting has:
- `book` — `score` (visible points) or `allowance` (daily spend budget)
- `amount` — signed
- `op_type` — `transfer` | `transfer_all` | `bet_win` | `bet_lose` | `admin_adjust` | `migration`
- `counterparty` — the *other* participant's telegram_id, or `0` for system ops (bets, admin)
- `op_id` — groups the postings of one logical operation

Domain enums live in `internal/domain/posting.go`. Time rules (UTC, midnight reset, "since" windows)
live in `internal/service/ledger.go` — reuse those helpers, don't reinvent `time.Now()`.

## Where things go

| You want to… | Touch |
|---|---|
| Add/change a user-facing string | `internal/i18n/i18n.go` (field) + `internal/i18n/ru.go` (value) |
| Add a command | handler in `internal/handlers/`, register in `handlers.go` |
| Add business logic | `internal/service/` |
| Add a DB query | `db/queries/*.sql` (+ `db/schema.sql` if new columns) → `task generate` → repository |
| Add a DB column/table | new migration (`task migrate:create`) **and** `db/schema.sql` (sqlc reads schema.sql) |
| Add a domain type/error | `internal/domain/` |
| Add a config/env var | `internal/config/config.go` + `.env.example` |

## Conventions & rules

- **sqlc output is generated and gitignored** (`/internal/database/gen/`). After editing
  `db/schema.sql` or `db/queries/*.sql`, run `task generate`. Never hand-edit `internal/database/gen`.
- **Adding a query** (pattern — copy an existing one like `Leaderboard` / `BetStats`):
  1. write the query in `db/queries/*.sql`
  2. `task generate`
  3. add a method to the repository interface + impl (`internal/repository/`), mapping gen rows → domain
  4. expose via a service (`internal/service/`)
  5. call from a handler
- **Adding a command** (pattern — see `internal/handlers/start.go` for the minimal shape):
  1. add the message field to `Messages` (`i18n.go`) and its value to `ru.go`
  2. write `func (h *Handlers) X(c tele.Context) error { defer logCommand(c, "/x")(); ... }`
  3. register it in `internal/handlers/handlers.go` `Register()`
  4. all user text via `h.msg.*` — never inline string literals in handlers
  - For wording, use the **i18n-tone** skill.
- **i18n voice**: lowercase, blunt, meme/rude, Russian. The term *несправление* is central
  ("не справился" = failed). Don't write polite/corporate copy. See `i18n-tone` skill.
- **Time is UTC.** Daily limits reset at midnight UTC. Use helpers in `internal/service/ledger.go`.
- **SQLite is single-connection** (`SetMaxOpenConns(1)` in `internal/database/database.go`) to avoid
  "database is locked". Migrations are embedded (`db/migrations/migrations.go`) and auto-applied on start.
  Create a migration with `task migrate:create -- <name>`.
- **Middleware order** (`internal/bot/bot.go` + `handlers.Register`): `recover` → rate limit →
  `ensureRegistered` → handler. Recover keeps one bad update from killing the bot.
- **Auto-registration**: the `ensureRegistered` middleware registers the sender on *any* message — there
  is no `/register`. New columns on `users` mean refreshing this path, not adding a command.
- **`ponytail:` comments** mark deliberate shortcuts with their upgrade path. Keep that convention.

## Build / dev / test

```bash
task tools          # install sqlc + goose (once)
task generate       # regenerate sqlc code (after editing db/)
task run            # run the bot (needs .env)
task test           # go test ./...
task build          # build bin/bot
task                # fmt + generate + build (pre-commit check)
task migrate        # apply migrations to $DB_PATH without starting the bot
task migrate:create -- <name>   # scaffold a new SQL migration
```

**Always run `task generate` before building if you changed anything under `db/`** — gen is not committed,
so a fresh checkout / CI relies on it being regenerated.

## Testing

- **Integration** (`internal/integration/`) — the `env` harness (`harness_test.go`) boots the real stack
  on a temp SQLite DB and exposes action+assertion helpers (`register`, `mustTransfer`, `score`,
  `remaining`, …). Tests read as scenarios, not plumbing. Prefer adding to this for end-to-end behavior.
- **Service unit tests** (e.g. `internal/service/user_test.go`) — small, table-ish, no frameworks.
- **No mocks.** Wire real repositories against a temp DB.

## Config

Env vars (see `.env.example`): `BOT_TOKEN` (required), `SUPER_ADMIN`, `DAILY_LIMIT` (default 1000),
`DB_PATH`, `DISCORD_TOKEN` (optional — enables Discord integration), `RATE_LIMIT_PER_SEC`,
`RATE_LIMIT_BURST`.

## Deployment

```bash
docker compose up -d --build
```

DB persists in `./data/`. The Docker build regenerates sqlc code and migrations run on startup.
