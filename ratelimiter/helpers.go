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
func NewLimiter(dispatcher *ext.Dispatcher, channels, edits bool) *Limiter {
	l := new(Limiter)

	l.filter = l.limiterFilter
	l.handler = l.limiterHandler
	l.timeout = DEFAULT_TIME
	l.punishment = DEFAULT_PUNISHMENT
	l.maxCount = DEFAULT_COUNT
	l.maxTimeout = DEFAULT_MAX_TIMEOUT
	l.IgnoreMediaGroup = true

	h := handlers.NewMessage(l.filter, l.handler)

	l.msgHandler = &h
	l.msgHandler.AllowChannel = channels
	l.msgHandler.AllowEdited = edits

	dispatcher.AddHandler(*l.msgHandler)
	return l
}

// NewFullLimiter creates a new `Limiter` with the given dispatcher.
// it will initialize a limiter which checks for messages received from
// channels and edited messages.
func NewFullLimiter(dispatcher *ext.Dispatcher) *Limiter {
	return NewLimiter(dispatcher, true, true)
}
