---
name: i18n-tone
description: Voice guide for the Несправлятор bot's user-facing strings. Use whenever writing or rewording any message in internal/i18n/ru.go (or any h.msg.* text). Ensures new copy matches the bot's aggressive lowercase Russian meme tone instead of sounding polite/corporate.
---

# Bot voice (i18n tone)

All user-facing text lives in `internal/i18n/ru.go` (`getRU()`), keyed by fields in
`internal/i18n/i18n.go`. The bot has a strong, consistent personality. Match it.

## The voice

- **Russian, lowercase.** No capitalized sentence starts. Emoji are welcome but not mandatory.
- **Blunt, rude, meme-y, playful.** It roasts the user. Mild profanity / chan-speak is on-brand
  ("педики", "лохи", "сосал??", "🥀"). It's a joke bot among friends, not a customer-service bot.
- **Built around "несправление".** The whole concept: *не справился* = you failed / screwed up.
  Points are "очки несправления". Lean on this vocabulary: справился / не справился / несправление.
- **Short.** One punchy line beats a paragraph. Errors are dismissive, not apologetic.

## Do / don't

- ✅ "ты уже играл сегодня, жадина. завтра приходи"
- ❌ "Вы уже использовали свою попытку на сегодня. Попробуйте завтра."
- ✅ "не справляюсь" (generic error)
- ❌ "Произошла ошибка. Пожалуйста, попробуйте позже."
- Keep `%` format verbs intact (`%s`, `%d`) and the same count/order when rewording.
- Don't translate to English. Don't go corporate, don't add please/sorry.

## Real examples (from ru.go)

| Field | Value |
|---|---|
| `MeError` / generic errors | `не справляюсь` |
| `BetAlreadyUsed` | `ты уже играл сегодня, жадина. завтра приходи` |
| `BetLose` | `-%d очков несправления 🥀🥀🥀` |
| `Help` (opening) | `📖 ликбез для даунов` |
| `StatsEmpty` | `пусто. напишите что-нибудь в чат чтобы встать на учёт, лохи` |
| `BalancesBetHint` | `🎲 помечены те, у кого не прокручена бетка, крутите педики` |
| `ReplySuccessNegative` | `%s не справился: %d` |
| `RateLimitError` | `сосал??` |

When adding a stat/feature line, name it literally but with attitude — e.g. the `/info` lines:
`🖖 навалил тебе минусов больше всех: %s (−%s)`, `🎯 по-твоему самый несправляющийся: %s (−%s)`.

If unsure between two phrasings, offer the user a few variants rather than guessing — the wording is
half the product here.
