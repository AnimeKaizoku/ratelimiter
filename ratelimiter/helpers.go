package ratelimiter

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func NewLimiter(dispatcher *ext.Dispatcher,
	channels, edits bool, tm ...int) *Limiter {
	l := new(Limiter)

	l.filter = l.limiterFilter
	l.handler = l.limiterHandler
	l.timeout = DEFAULT_TIME
	l.punishment = DEFAULT_PUNISHMENT
	l.maxCount = DEFAULT_COUNT
	l.maxTimeout = DEFAULT_MAX_TIMEOUT
	l.IgnoreMediaGroup = true

	msgHandler := handlers.NewMessage(l.filter, l.handler)
	msgHandler.AllowChannel = channels
	msgHandler.AllowEdited = edits

	dispatcher.AddHandler(msgHandler)
	return l
}

func NewFullLimiter(dispatcher *ext.Dispatcher) *Limiter {
	return NewLimiter(dispatcher, true, true)
}
