---
name: add-command
description: Add a new Telegram command to the Несправлятор bot end-to-end. Use when asked to "add a command", "new bot command", "add a /xxx handler", or wire up a new slash command. Covers the i18n field, handler file, registration, and conventions.
---

# Add a Telegram command

The bot routes commands through telebot in `internal/handlers`. Adding one is a fixed 4-step pattern.
Keep handlers thin: parse input, call a service, send a localized string. No business logic, no inline text.

## Steps

1. **i18n field + value.** Add the message field to the `Messages` struct in
   `internal/i18n/i18n.go`, then its value in `internal/i18n/ru.go` (`getRU()`).
   All user-facing text lives here — never inline a string literal in a handler.
   For the wording/voice, use the **i18n-tone** skill.

2. **Handler.** Create `internal/handlers/<name>.go`. Minimal shape (mirror `start.go`):

   ```go
   package handlers

   import tele "gopkg.in/telebot.v3"

   func (h *Handlers) Foo(c tele.Context) error {
       defer logCommand(c, "/foo")()
       // parse c.Message(), call a service (h.userService / h.balanceService / h.betService / ...)
       return c.Send(h.msg.FooSomething)
   }
   ```

   - Always `defer logCommand(c, "/foo")()` first — it logs req_id + duration.
   - Reach services through the `h.*Service` fields; reach repos only if a service can't fit.
   - Reply-target logic: `c.Message().ReplyTo` (see `info.go` / `reply.go` for the pattern,
     including the "replying to the bot" case via `h.bot.Me.ID`).

3. **Register.** Add it in `internal/handlers/handlers.go` `Register()`:
   `h.bot.Handle("/foo", h.Foo)`.

4. **Build & verify.** `task build` (or `go build ./...`). If the command needs new data, add a query
   first — but that's a separate flow (see CLAUDE.md "Adding a query").

## Notes

- Senders are auto-registered by the `ensureRegistered` middleware on any message, so you don't need a
  registration/`GetOrCreate` call in the handler.
- Helpers already in the package: `displayName(username, firstName)` and the number/trend formatters in
  `info.go`. Reuse them instead of re-implementing.
- If the command is admin-only, follow `add.go` (super-admin check via the balance service).
- Consider whether `/help` (`ru.go`) should mention the new command.
