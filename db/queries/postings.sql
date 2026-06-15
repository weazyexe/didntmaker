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

-- name: GetUserBetStats :one
SELECT
    CAST(COUNT(CASE WHEN op_type = 'bet_win'  THEN 1 END) AS INTEGER) AS won,
    CAST(COUNT(CASE WHEN op_type = 'bet_lose' THEN 1 END) AS INTEGER) AS lost
FROM postings
WHERE chat_id = ? AND account_id = ? AND op_type IN ('bet_win', 'bet_lose');

-- name: GetChatBetStats :many
SELECT
    u.telegram_id,
    u.username,
    u.first_name,
    CAST(COUNT(CASE WHEN p.op_type = 'bet_win'  THEN 1 END) AS INTEGER) AS won,
    CAST(COUNT(CASE WHEN p.op_type = 'bet_lose' THEN 1 END) AS INTEGER) AS lost
FROM users u
JOIN postings p
    ON p.chat_id = u.chat_id
   AND p.account_id = u.telegram_id
   AND p.op_type IN ('bet_win', 'bet_lose')
WHERE u.chat_id = ?
GROUP BY u.id
ORDER BY (won - lost) DESC;

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

-- name: GetChatBetAccountsSince :many
SELECT DISTINCT account_id FROM postings
WHERE chat_id = ? AND op_type IN ('bet_win', 'bet_lose') AND created_at >= ?;

-- name: GetScoreSince :one
SELECT CAST(COALESCE(SUM(amount), 0) AS INTEGER) AS delta
FROM postings
WHERE chat_id = ? AND account_id = ? AND book = 'score' AND created_at >= ?;

-- name: GetIncomingByCounterparty :many
SELECT
    u.username,
    u.first_name,
    CAST(COALESCE(SUM(CASE WHEN p.amount > 0 THEN p.amount ELSE 0 END), 0) AS INTEGER) AS plus,
    CAST(COALESCE(SUM(CASE WHEN p.amount < 0 THEN -p.amount ELSE 0 END), 0) AS INTEGER) AS minus
FROM postings p
JOIN users u
    ON u.chat_id = p.chat_id
   AND u.telegram_id = p.counterparty
WHERE p.chat_id = ? AND p.account_id = ? AND p.book = 'score' AND p.counterparty != 0
GROUP BY p.counterparty;

-- name: GetOutgoingByAccount :many
SELECT
    u.username,
    u.first_name,
    CAST(COALESCE(SUM(CASE WHEN p.amount > 0 THEN p.amount ELSE 0 END), 0) AS INTEGER) AS plus,
    CAST(COALESCE(SUM(CASE WHEN p.amount < 0 THEN -p.amount ELSE 0 END), 0) AS INTEGER) AS minus
FROM postings p
JOIN users u
    ON u.chat_id = p.chat_id
   AND u.telegram_id = p.account_id
WHERE p.chat_id = ? AND p.counterparty = ? AND p.book = 'score'
GROUP BY p.account_id;
