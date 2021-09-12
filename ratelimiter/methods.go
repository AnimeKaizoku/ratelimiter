package ratelimiter

import (
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

//---------------------------------------------------------

func (l *Limiter) Start() {
	if l.isEnabled {
		return
	}

	if l.mutex == nil {
		l.mutex = new(sync.Mutex)
	}

	if l.userMap == nil {
		l.userMap = make(map[int64]*UserStatus)
	}

	go l.checker()

	l.isEnabled = true
	l.isStopped = false
}

func (l *Limiter) Stop() {
	if l.isStopped {
		return
	}

	l.isEnabled = false
	l.isStopped = true
}

func (l *Limiter) IsStopped() bool {
	return l.isStopped
}

func (l *Limiter) IsEnabled() bool {
	return l.isEnabled
}

// SetTriggerFunc will set the trigger function of this limiter.
// The trigger function will be triggered when the limiter
// limits a user. The information passed by it will be the
// information related to the last message of the user.
func (l *Limiter) SetTriggerFunc(t handlers.Response) {
	l.trigger = t
}

func (l *Limiter) AddException(ex filters.Message) {
	l.exceptions = append(l.exceptions, ex)
}

func (l *Limiter) ClearAllExceptions() {
	l.exceptions = nil
}

func (l *Limiter) GetExceptions() []filters.Message {
	return l.exceptions
}

func (l *Limiter) AddExceptionID(id int64) {
	l.exceptionIDs = append(l.exceptionIDs, id)
}

func (l *Limiter) ClearAllExceptionIDs() {
	l.exceptionIDs = nil
}

func (l *Limiter) IsInExcpetionList(id int64) bool {
	if len(l.exceptionIDs) == 0 {
		return false
	}

	for _, ex := range l.exceptionIDs {
		if ex == id {
			return true
		}
	}

	return false
}

func (l *Limiter) isException(msg *gotgbot.Message) bool {
	if len(l.exceptionIDs) == 0 {
		return false
	}

	for _, ex := range l.exceptionIDs {
		if msg.From != nil {
			if ex == msg.From.Id || ex == msg.Chat.Id {
				return true
			}
		} else {
			if ex == msg.Chat.Id {
				return true
			}
		}

	}

	return false
}

func (l *Limiter) checker() {
	for l.isEnabled && !l.isStopped {
		time.Sleep(l.maxTimeout)
		if len(l.userMap) == 0 {
			continue
		}

		l.mutex.Lock()
		for key, value := range l.userMap {
			if time.Since(value.Last) > l.timeout {
				delete(l.userMap, key)
			}
		}
		l.mutex.Unlock()
	}
}

//---------------------------------------------------------
