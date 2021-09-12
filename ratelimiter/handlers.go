package ratelimiter

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (l *Limiter) limiterFilter(msg *gotgbot.Message) bool {
	if !l.isEnabled || l.isStopped {
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

	return !(l.IgnoreMediaGroup && msg.MediaGroupId != "")
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
		if l.trigger != nil {
			go l.trigger(b, ctx)
		}

		return ext.EndGroups
	}

	status.Last = time.Now()
	l.mutex.Unlock()

	return ext.ContinueGroups
}
