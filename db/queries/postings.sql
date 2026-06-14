-- name: InsertPosting :exec
INSERT INTO postings (
    chat_id, account_id, book, amount, op_id, op_type, counterparty, metadata, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetScore :one
SELECT CAST(COALESCE(SUM(amount), 0) AS INTEGER) AS score
FROM postings
WHERE chat_id = ? AND account_id = ? AND book = 'score';

-- name: GetAllowanceSpentSince :one
SELECT CAST(COALESCE(SUM(amount), 0) AS INTEGER) AS spent
FROM postings
WHERE chat_id = ? AND account_id = ? AND book = 'allowance' AND created_at >= ?;

-- name: GetChatAllowanceSpentSince :many
SELECT account_id, CAST(COALESCE(SUM(amount), 0) AS INTEGER) AS spent
FROM postings
WHERE chat_id = ? AND book = 'allowance' AND created_at >= ?
GROUP BY account_id;

-- name: HasBetSince :one
SELECT EXISTS (
    SELECT 1 FROM postings
    WHERE chat_id = ? AND account_id = ?
      AND op_type IN ('bet_win', 'bet_lose')
      AND created_at >= ?
) AS has_bet;

-- name: GetLeaderboard :many
SELECT
    u.telegram_id,
    u.username,
    u.first_name,
    CAST(COALESCE(SUM(p.amount), 0) AS INTEGER) AS score
FROM users u
LEFT JOIN postings p
    ON p.chat_id = u.chat_id
   AND p.account_id = u.telegram_id
   AND p.book = 'score'
WHERE u.chat_id = ?
GROUP BY u.id
ORDER BY score ASC;
