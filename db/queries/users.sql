-- name: GetUserByTelegramID :one
SELECT * FROM users
WHERE chat_id = ? AND telegram_id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE chat_id = ? AND username = ?;

-- name: CreateUser :one
INSERT INTO users (telegram_id, chat_id, username, first_name, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUserProfile :exec
UPDATE users
SET username = ?, first_name = ?, updated_at = ?
WHERE chat_id = ? AND telegram_id = ?;

-- name: ListChatUsers :many
SELECT * FROM users
WHERE chat_id = ?;

-- name: CountChatUsersExcept :one
SELECT COUNT(*) FROM users
WHERE chat_id = ? AND telegram_id != ?;

-- name: ListChatUserIDsExcept :many
SELECT telegram_id FROM users
WHERE chat_id = ? AND telegram_id != ?;
