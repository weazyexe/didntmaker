-- name: CreateDiscordBinding :exec
INSERT INTO discord_bindings (chat_id, guild_id, created_at)
VALUES (?, ?, ?);

-- name: DeleteDiscordBinding :exec
DELETE FROM discord_bindings
WHERE chat_id = ? AND guild_id = ?;

-- name: GetDiscordBindingsByGuildID :many
SELECT * FROM discord_bindings
WHERE guild_id = ?;

-- name: GetDiscordBindingsByChatID :many
SELECT * FROM discord_bindings
WHERE chat_id = ?;

-- name: DiscordBindingExists :one
SELECT EXISTS (
    SELECT 1 FROM discord_bindings
    WHERE chat_id = ? AND guild_id = ?
) AS exists_binding;
