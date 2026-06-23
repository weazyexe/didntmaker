package i18n

// Lang represents a language code
type Lang string

const (
	LangRU Lang = "ru"
)

// Messages contains all bot messages
type Messages struct {
	// Commands
	Start string
	Help  string

	// /info command (stats layout)
	MeHeader   string
	MeScore    string
	MeWeek     string
	MeMonth    string
	MeWorstDay string
	MeBestDay  string
	MeLimit    string
	MeBets     string
	MeBetReady string
	MeBetUsed  string
	MeFan      string
	MeHater    string
	MeFavorite string
	MeVictim   string
	MeError    string

	// /info command
	InfoOnBot string

	// /bet_stats command
	BetStatsHeader string
	BetStatsEmpty  string
	BetStatsError  string
	BetStatsEntry  string

	// /stats command
	StatsHeader      string
	StatsWeekHeader  string
	StatsMonthHeader string
	StatsEmpty       string
	StatsError       string
	StatsMedals      []string
	StatsEntryFmt    string

	// /balances command
	BalancesHeader       string
	BalancesEmpty        string
	BalancesError        string
	BalancesEntry        string
	BalancesBetAvailable string
	BalancesBetHint      string

	// /bet command
	BetAlreadyUsed string
	BetError       string
	BetDiceError   string
	BetResultError string
	BetWin         string
	BetLose        string

	// /add command (admin)
	AddUsage       string
	AddFormatError string
	AddNumberError string
	AddNotFound    string
	AddError       string
	AddSuccess     string

	// Reply handler
	ReplyUnknownTarget   string
	ReplySelfError       string
	ReplyNotEnough       string
	ReplyTargetNotFound  string
	ReplyError           string
	ReplySuccessNegative string
	ReplySuccessPositive string
	ReplyAllNoUsers      string
	ReplyAllNotEnough    string
	ReplyAllError        string
	ReplyAllSuccessNeg   string
	ReplyAllSuccessPos   string

	// Discord voice events
	DiscordVoiceJoin  string
	DiscordVoiceLeave string

	// Discord bind commands
	DiscordBindUsage        string
	DiscordBindInvalidID    string
	DiscordBindAlreadyBound string
	DiscordBindError        string
	DiscordBindSuccess      string
	DiscordUnbindUsage      string
	DiscordUnbindError      string
	DiscordUnbindSuccess    string

	// service strings
	RateLimitError string
}

// Default returns the default language (Russian)
func Default() Lang {
	return LangRU
}

// Get returns messages for the specified language
func Get(lang Lang) *Messages {
	switch lang {
	case LangRU:
		fallthrough
	default:
		return getRU()
	}
}
