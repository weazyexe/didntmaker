-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id    INTEGER  NOT NULL,
    chat_id        INTEGER  NOT NULL,
    username       TEXT     NOT NULL DEFAULT '',
    first_name     TEXT     NOT NULL DEFAULT '',
    balance        INTEGER  NOT NULL DEFAULT 0,
    daily_given    INTEGER  NOT NULL DEFAULT 0,
    daily_reset_at DATETIME,
    last_bet_at    DATETIME,
    created_at     DATETIME NOT NULL,
    updated_at     DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_telegram_chat ON users (telegram_id, chat_id);

CREATE TABLE IF NOT EXISTS discord_bindings (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    chat_id    INTEGER  NOT NULL,
    guild_id   TEXT     NOT NULL,
    created_at DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_discord_binding ON discord_bindings (chat_id, guild_id);
CREATE INDEX IF NOT EXISTS idx_guild ON discord_bindings (guild_id);

-- +goose Down
DROP TABLE IF EXISTS discord_bindings;
DROP TABLE IF EXISTS users;
