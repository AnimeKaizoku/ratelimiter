package ratelimiter

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func NewLimiter(dispatcher *ext.Dispatcher) *Limiter {
	l := new(Limiter)

	l.filter = l.limiterFilter
	l.handler = l.limiterHandler
	l.timeout = DEFAULT_TIME
	l.maxCount = DEFAULT_COUNT
	l.maxTimeout = DEFAULT_MAX_TIMEOUT
	l.IgnoreMediaGroup = true

	msgHandle := handlers.NewMessage(l.filter, l.handler)

	dispatcher.AddHandler(msgHandle)
	return l
}
