package i18n

func getRU() *Messages {
	return &Messages{
		// /start
		Start: `это бот несправления. чем меньше очков — тем больше ты не справился по жизни.

реплай на сообщение с +N или -N — записать несправление
/info — глянуть свой позор (или чужой, если реплайнуть)
/stats — топ несправляющихся
/bet — проиграй ещё больше
/help — если совсем тупой`,

		// /help
		Help: `📖 ликбез для даунов

🎯 зачем этот бот
считает очки несправления. не справился — получи минус. справился — держи плюс (но это не точно).

⚡ как пользоваться

реплай +N или -N
реплаишь на сообщение человека и пишешь число. всё.
-100 значит не справился на сотку
+50 значит немного справился (чё, бывает)

реплай на бота +N или -N
реплаишь на любое сообщение бота — баллы летят всем кроме тебя.
стоимость = очки × количество людей в чате

/info
посмотреть насколько ты в жопе и сколько ещё можешь раздать другим. реплайни на чьё-то сообщение с /info — глянешь его позор
учёт автоматический: просто напиши что-нибудь в чат и ты в системе

/balances
у кого сколько очков несправления на раздачу осталось на сегодня

/stats
рейтинг несправляющихся. лёша не справляется по умолчанию

/bet
кинь кубик, тряхни очком 🎲
выпало 4-6 → получаешь +1000 к лимиту раздачи
выпало 1-3 → сам получаешь 1000 очков несправления
один раз в день, для самых отчаянных пидоров

📏 правила

• в день можешь раздать 1000 очков (и еще 1000, если выиграл бетку, петух)
• себе накинуть нельзя
• в плюс шанс выйти мал, ну ты попробуй, малой
• лимиты сбрасываются в полночь по UTC

🔊 discord
бот умеет кидать уведомления когда кто-то заходит/выходит из войса на discord сервере

/discord_bind <guild_id> — привязать discord сервер к этому чату
/discord_unbind <guild_id> — отвязать

guild_id — это ID сервера в discord. включи Developer Mode (настройки → расширенные), потом ПКМ по серверу → Copy Server ID

💡 пример
[реплай на сообщение леши]
-1000
→ леша не справился: -1000`,

		// /info
		MeHeader:   "👤 %s\n\n",
		MeScore:    "🏆 очки несправления: %s\n",
		MeWeek:     "📈 за неделю: %s\n",
		MeMonth:    "📅 за месяц: %s\n",
		MeLimit:    "\n🎟 раздать сегодня: %s из %s\n",
		MeBets:     "🎲 бетки: %d справился / %d не справился",
		MeBetReady: " · кубик готов 🎲\n",
		MeBetUsed:  " · кубик уже крутил\n",
		MeFan:      "❤️ дал тебе больше всего плюсов: %s (+%s)\n",
		MeHater:    "🖤 навалил тебе минусов больше всех: %s (−%s)\n",
		MeFavorite: "😎 ты дал плюсов больше всего: %s (+%s)\n",
		MeVictim:   "🎯 по-твоему самый несправляющийся: %s (−%s)\n",
		MeError:    "не справляюсь",

		// /info
		InfoOnBot: "сосал??",

		// /bet_stats
		BetStatsHeader: "🎲 статистика бетов:\n\n",
		BetStatsEmpty:  "ещё никто не крутил бетку",
		BetStatsError:  "не справляюсь",
		BetStatsEntry:  "%s: выиграл раз - %d, проиграл раз - %d\n",

		// /stats
		StatsHeader:      "🏆 топ несправляющихся:\n\n",
		StatsWeekHeader:  "🏆 топ несправляющихся за неделю:\n\n",
		StatsMonthHeader: "🏆 топ несправляющихся за месяц:\n\n",
		StatsEmpty:       "пусто. напишите что-нибудь в чат чтобы встать на учёт, лохи",
		StatsError:       "не могу достать статистику, не справляюсь",
		StatsMedals:      []string{"🥇", "🥈", "🥉"},
		StatsEntryFmt:    "%s %s: %d\n",

		// /balances
		BalancesHeader:       "🔫 лимит очков несправления на раздачу на сегодня:\n\n",
		BalancesEmpty:        "никого нет. напишите что-нибудь чтобы зарегаться",
		BalancesError:        "не справляюсь",
		BalancesEntry:        "%s: %d/%d",
		BalancesBetAvailable: " 🎲",
		BalancesBetHint:      "🎲 помечены те, у кого не прокручена бетка, крутите педики",

		// /bet
		BetAlreadyUsed: "ты уже играл сегодня, жадина. завтра приходи",
		BetError:       "не справляюсь",
		BetDiceError:   "кубик сломался, не справляюсь",
		BetResultError: "не смог записать результат, не справляюсь",
		BetWin:         "хоть когда-то ты справился, +%d к лимиту",
		BetLose:        "-%d очков несправления 🥀🥀🥀",

		// /add (admin)
		AddUsage:       "/add @username +/-N",
		AddFormatError: "формат: /add @username +/-N",
		AddNumberError: "число введи нормально",
		AddNotFound:    "@%s не найден",
		AddError:       "не справляюсь",
		AddSuccess:     "право сильнейшего: @%s: %d → %d (%s%d)",

		// Reply handler
		ReplyUnknownTarget:   "сосал??",
		ReplySelfError:       "ты не справился",
		ReplyNotEnough:       "не хватает. осталось %d из %d",
		ReplyTargetNotFound:  "этот чел не зарегался. пусть напишет что-нибудь в чат",
		ReplyError:           "не справляюсь",
		ReplySuccessNegative: "%s не справился: %d",
		ReplySuccessPositive: "%s справился: +%d",

		// Reply to all
		ReplyAllNoUsers:    "ты одинокий педик, раздавать некому",
		ReplyAllNotEnough:  "не хватает. надо %d, осталось %d",
		ReplyAllError:      "не справляюсь",
		ReplyAllSuccessNeg: "все не справились: %d каждому",
		ReplyAllSuccessPos: "все справились: +%d каждому",

		// Discord voice events
		DiscordVoiceJoin:  "%s присоединился к каналу #%s",
		DiscordVoiceLeave: "%s покинул канал #%s",

		// Discord bind commands
		DiscordBindUsage: `формат: /discord_bind <guild_id>

guild_id — это ID сервера в discord. чтобы получить:
1. открой настройки discord → расширенные → включи Developer Mode
2. ПКМ по серверу → Copy Server ID
3. /discord_bind <скопированный_id>`,
		DiscordBindInvalidID:    "guild_id должен быть числом",
		DiscordBindAlreadyBound: "этот сервер уже привязан к этому чату",
		DiscordBindError:        "не справляюсь",
		DiscordBindSuccess:      "discord сервер %s привязан к этому чату",
		DiscordUnbindUsage:      "формат: /discord_unbind <guild_id>",
		DiscordUnbindError:      "не справляюсь",
		DiscordUnbindSuccess:    "discord сервер %s отвязан от этого чата",

		RateLimitError: "сосал??",
	}
}
