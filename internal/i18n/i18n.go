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

	// /me command
	MeStats string
	MeError string

	// /stats command
	StatsHeader   string
	StatsEmpty    string
	StatsError    string
	StatsMedals   []string
	StatsEntryFmt string

	// /balances command
	BalancesHeader string
	BalancesEmpty  string
	BalancesError  string
	BalancesFull   string
	BalancesEmpty_ string
	BalancesEntry  string

	// /bet command
	BetNotRegistered string
	BetAlreadyUsed   string
	BetError         string
	BetDiceError     string
	BetResultError   string
	BetWin           string
	BetLose          string

	// /add command (admin)
	AddUsage       string
	AddFormatError string
	AddNumberError string
	AddNotFound    string
	AddError       string
	AddSuccess     string

	// Reply handler
	ReplyLimitExceeded   string
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
	ReplyNotRegistered   string
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
