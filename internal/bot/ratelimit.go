package bot

import (
	"sync"

	"golang.org/x/time/rate"
	tele "gopkg.in/telebot.v3"
	"weazyexe.dev/didntmaker/internal/i18n"
)

type rateLimiter struct {
	mu    sync.Mutex
	r     rate.Limit
	burst int
	users map[int64]*rate.Limiter
	msg   *i18n.Messages
}

func newRateLimiter(r rate.Limit, burst int, msg *i18n.Messages) *rateLimiter {
	return &rateLimiter{r: r, burst: burst, users: make(map[int64]*rate.Limiter), msg: msg}
}

func (rl *rateLimiter) allow(userID int64) bool {
	rl.mu.Lock()
	l, ok := rl.users[userID]
	if !ok {
		l = rate.NewLimiter(rl.r, rl.burst)
		rl.users[userID] = l
	}
	rl.mu.Unlock()
	return l.Allow()
}

// Middleware silently drops updates from users over their rate.
func (rl *rateLimiter) Middleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if s := c.Sender(); s != nil && !rl.allow(s.ID) {
				c.Send(rl.msg.RateLimitError)
				return nil
			}
			return next(c)
		}
	}
}
