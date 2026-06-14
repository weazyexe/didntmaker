-- Canonical schema used by sqlc for code generation.
-- The runtime schema is produced by the goose migrations in db/migrations
-- and must stay in sync with this file.

CREATE TABLE users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER  NOT NULL,
    chat_id     INTEGER  NOT NULL,
    username    TEXT     NOT NULL DEFAULT '',
    first_name  TEXT     NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL
);
CREATE UNIQUE INDEX idx_telegram_chat ON users (telegram_id, chat_id);

-- Append-only ledger. Source of truth for both score and daily allowance.
-- One user action = several postings sharing the same op_id.
CREATE TABLE postings (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    chat_id      INTEGER  NOT NULL,
    account_id   INTEGER  NOT NULL,            -- telegram_id of the affected account
    book         TEXT     NOT NULL,            -- 'score' | 'allowance'
    amount       INTEGER  NOT NULL,            -- signed
    op_id        TEXT     NOT NULL,            -- groups postings of one operation
    op_type      TEXT     NOT NULL,            -- transfer|transfer_all|bet_win|bet_lose|admin_adjust|migration
    counterparty INTEGER  NOT NULL DEFAULT 0,  -- the other participant, for stats
    metadata     TEXT     NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL
);
CREATE INDEX idx_score ON postings (chat_id, account_id, book);
CREATE INDEX idx_allow ON postings (chat_id, account_id, book, created_at);
CREATE INDEX idx_counter ON postings (chat_id, account_id, book, created_at, counterparty);

CREATE TABLE discord_bindings (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    chat_id    INTEGER  NOT NULL,
    guild_id   TEXT     NOT NULL,
    created_at DATETIME NOT NULL
);
CREATE UNIQUE INDEX idx_discord_binding ON discord_bindings (chat_id, guild_id);
CREATE INDEX idx_guild ON discord_bindings (guild_id);
