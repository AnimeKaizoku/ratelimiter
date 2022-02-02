// ratelimiter Project
// Copyright (C) 2021 ALiwoto and other Contributors
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ratelimiter

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// NewLimiter creates a new `Limiter` with the given dispatcher.
// pass true for the second parameter if you want the limiter to check
// messages in channels too.
// pass true for the third parameter if you want the limiter to check
// edited messages too.
func NewLimiter(dispatcher *ext.Dispatcher, config *LimiterConfig) *Limiter {
	l := new(Limiter)

	if config == nil {
		config = &LimiterConfig{
			ConsiderChannel:  false,
			ConsiderUser:     true,
			ConsiderEdits:    false,
			IgnoreMediaGroup: true,
			TextOnly:         false,
			ConsiderInline:   true,
			IsStrict:         false,
		}
	}

	l.filter = l.limiterFilter
	l.handler = l.limiterHandler
	l.timeout = DEFAULT_TIME
	l.punishment = DEFAULT_PUNISHMENT
	l.maxCount = DEFAULT_COUNT
	l.maxTimeout = DEFAULT_MAX_TIMEOUT
	l.IgnoreMediaGroup = config.IgnoreMediaGroup
	l.TextOnly = config.TextOnly
	l.ConsiderUser = config.ConsiderUser
	l.ConsiderInline = config.ConsiderInline
	l.IsStrict = config.IsStrict

	h := handlers.NewMessage(l.filter, l.handler)
	cb := handlers.NewCallback(l.callbackFilter, l.handler)

	l.msgHandler = &h
	l.msgHandler.AllowChannel = config.ConsiderChannel
	l.msgHandler.AllowEdited = config.ConsiderEdits

	l.allHandlers = append(l.allHandlers, h, cb)

	for _, currentHandler := range l.allHandlers {
		if len(config.HandlerGroups) != 0 {
			for _, current := range config.HandlerGroups {
				dispatcher.AddHandlerToGroup(currentHandler, current)
			}
		} else {
			dispatcher.AddHandler(currentHandler)
		}
	}

	return l
}

// NewFullLimiter creates a new `Limiter` with the given dispatcher.
// it will initialize a limiter which checks for messages received from
// channels and edited messages.
func NewFullLimiter(dispatcher *ext.Dispatcher) *Limiter {
	return NewLimiter(dispatcher, &LimiterConfig{
		ConsiderChannel:  true,
		ConsiderUser:     true,
		ConsiderEdits:    true,
		IgnoreMediaGroup: false,
		TextOnly:         false,
		IsStrict:         false,
		ConsiderInline:   true,
	})
}
