-- +goose Up
CREATE TABLE IF NOT EXISTS postings (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    chat_id      INTEGER  NOT NULL,
    account_id   INTEGER  NOT NULL,
    book         TEXT     NOT NULL,
    amount       INTEGER  NOT NULL,
    op_id        TEXT     NOT NULL,
    op_type      TEXT     NOT NULL,
    counterparty INTEGER  NOT NULL DEFAULT 0,
    metadata     TEXT     NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_score ON postings (chat_id, account_id, book);
CREATE INDEX IF NOT EXISTS idx_allow ON postings (chat_id, account_id, book, created_at);
CREATE INDEX IF NOT EXISTS idx_counter ON postings (chat_id, account_id, book, created_at, counterparty);

INSERT INTO postings (chat_id, account_id, book, amount, op_id, op_type, counterparty, metadata, created_at)
SELECT chat_id, telegram_id, 'score', balance, 'migration-' || id, 'migration', 0, '', created_at
FROM users
WHERE balance != 0;

ALTER TABLE users DROP COLUMN balance;
ALTER TABLE users DROP COLUMN daily_given;
ALTER TABLE users DROP COLUMN daily_reset_at;
ALTER TABLE users DROP COLUMN last_bet_at;

DROP TABLE IF EXISTS transactions;

-- +goose Down
DROP TABLE IF EXISTS postings;
