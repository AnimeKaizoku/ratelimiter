// ratelimiter Project
// Copyright (C) 2021 ALiwoto and other Contributors
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ratelimiter

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (l *Limiter) limiterFilter(msg *gotgbot.Message) bool {
	if !l.isEnabled || l.isStopped || !l.hasTextCondition(msg) {
		return false
	}

	if l.isException(msg) {
		return false
	}

	if len(l.exceptions) != 0 {
		for _, ex := range l.exceptions {
			if ex(msg) {
				return false
			}
		}
	}

	if len(l.conditions) != 0 {
		for _, con := range l.conditions {
			if !con(msg) {
				return false
			}
		}
	}

	return !(l.IgnoreMediaGroup && len(msg.MediaGroupId) != 0)
}

func (l *Limiter) limiterHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	var status *UserStatus
	var id int64
	if l.ConsiderUser && ctx.EffectiveUser != nil {
		id = ctx.EffectiveUser.Id
	} else if ctx.EffectiveChat != nil {
		id = ctx.EffectiveChat.Id
	} else {
		return ext.ContinueGroups
	}

	l.mutex.Lock()
	status = l.userMap[id]
	if status == nil {
		status = new(UserStatus)
		status.Last = time.Now()
		status.count++
		l.userMap[id] = status
		l.mutex.Unlock()
		if status.IsCustomLimited() {
			if !status.custom.ignoreException && l.isException(ctx.Message) {
				return ext.ContinueGroups
			}
			return ext.EndGroups
		}
		return ext.ContinueGroups
	}

	if status.limited {
		l.mutex.Unlock()
		if time.Since(status.Last) > l.timeout+l.punishment {
			status.count = 0
			status.limited = false
			status.Last = time.Now()
			return ext.ContinueGroups
		}

		if l.IsStrict {
			status.Last = time.Now()
		}

		return ext.EndGroups
	}

	if time.Since(status.Last) > l.timeout {
		status.count = 0
	}

	status.count++

	if status.count > l.maxCount {
		status.limited = true
		status.Last = time.Now()
		l.mutex.Unlock()
		// check for triggers length to prevent from creating
		// a new goroutine in the case we have no triggers.
		if len(l.triggers) != 0 {
			go l.runTriggers(b, ctx)
		}

		return ext.EndGroups
	}

	l.mutex.Unlock()
	status.Last = time.Now()

	if status.IsCustomLimited() {
		if !status.custom.ignoreException && l.isException(ctx.Message) {
			return ext.ContinueGroups
		}
		return ext.EndGroups
	}

	return ext.ContinueGroups
}
