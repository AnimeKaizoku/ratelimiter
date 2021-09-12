package ratelimiter

import (
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

//---------------------------------------------------------

// Start will start the limiter.
// When the limiter is started (enabled), it will check for
// check for incoming messages; if they are considered as flood,
// the limiter won't let the handler functions to be called.
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

	l.isEnabled = true
	l.isStopped = false

	go l.checker()
}

// Stop method will make this limiter stop checking the incoming
// messages and will set its variables to nil.
// the main resources used by this limiter will be freed,
// such as map and mutex.
// but the configuration variables such as message time out will
// remain the same and won't be set to 0.
func (l *Limiter) Stop() {
	if l.isStopped {
		return
	}

	l.isEnabled = false
	l.isStopped = true

	// make sure that mutex is not nil.
	if l.mutex != nil {
		// let another goroutines let go of the mutex;
		// if you set the usermap value to nil out of nowhere
		// it MAY cause some troubles.
		l.mutex.Lock()
		l.userMap = nil
		l.mutex.Unlock()

		l.mutex = nil
	}
}

// IsStopped returns true if this limiter is already stopped
// and doesn't check for incoming messages.
func (l *Limiter) IsStopped() bool {
	return l.isStopped
}

// IsEnabled returns true if and only if this limiter is enabled
// and is checking the incoming messages for floodwait.
// for enabling the limiter, you need to use `Start` method.
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

// AddException will add an exception filter to this limiter.
func (l *Limiter) AddException(ex filters.Message) {
	l.exceptions = append(l.exceptions, ex)
}

// ClearAllExceptions will clear all exception of this limiter.
// this way, you will be sure that all of incoming updates will be
// checked for floodwait by this limiter.
func (l *Limiter) ClearAllExceptions() {
	l.exceptions = nil
}

// GetExceptions returns the filters array used by this limiter as
// its exceptions list.
func (l *Limiter) GetExceptions() []filters.Message {
	return l.exceptions
}

// AddExceptionID will add a group/user/channel ID to the exception
// list of the limiter.
func (l *Limiter) AddExceptionID(id ...int64) {
	l.exceptionIDs = append(l.exceptionIDs, id...)
}

// AddCondition will add a condition to be checked by this limiter,
// if this condition doesn't return true, the limiter won't check
// the message for antifloodwait.
func (l *Limiter) AddCondition(condition filters.Message) {
	l.conditions = append(l.conditions, condition)
}

// ClearAllConditions clears all condition list.
func (l *Limiter) ClearAllConditions() {
	l.conditions = nil
}

// AddConditions will accept an array of the conditions and will
// add them to the condition list of this limiter.
// you can also pass only one value to this method.
func (l *Limiter) AddConditions(conditions ...filters.Message) {
	l.conditions = append(l.conditions, conditions...)
}

// SetAsConditions will accept an array of conditions and will set
// the conditions of the limiter to them.
func (l *Limiter) SetAsConditions(conditions []filters.Message) {
	l.conditions = conditions
}

// ClearAllExceptions will clear all exception IDs of this limiter.
// this way, you will be sure that all of incoming updates will be
// checked for floodwait by this limiter.
func (l *Limiter) ClearAllExceptionIDs() {
	l.exceptionIDs = nil
}

// IsInExcpetionList will check and see if an ID is in the
// exception list of the listener or not.
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

// SetAsExceptionList will set its argument at the exception
// list of this limiter. Please notice that this method won't
// append the list to the already existing exception list; but
// it will set it to this, so the already existing exception IDs
// assigned to this limiter will be lost.
func (l *Limiter) SetAsExceptionList(list []int64) {
	l.exceptionIDs = list
}

// GetStatus will get the status of a chat.
// if `l.ConsiderUser` parameter is set to `true`,
// the id should be the id of the user; otherwise you should
// use the id of the chat to get the status.
func (l *Limiter) GetStatus(id int64) *UserStatus {
	var status *UserStatus
	l.mutex.Lock()
	status = l.userMap[id]
	l.mutex.Unlock()

	return status
}

// SetFloodWaitTime will set the flood wait duration for each
// chat to send `maxCount` message per this amount of time.
// if they send more than this amount of messages during this time,
// they will be limited by this limiter and so their messages
// won't be handled in the current group.
// (Notice: if `ConsiderUser` is set to `true`, this duration will
// be applied to unique users in the chat; not the total chat.)
func (l *Limiter) SetFloodWaitTime(d time.Duration) {
	l.timeout = d
}

// SetPunishmentDuration will set the punishment duration of
// the chat (or a user) after being limited by this limiter.
// Users needs to spend this amount of time + `l.timeout` to become
// free and so the handlers will work again for it.
// NOTICE: if `IsStrict` is set to `true`, as long as user sends
// message to the bot, the amount of passed-punishment time will
// become 0; so the user needs to stop sending messages to the bot
// until the punishment time is passed, otherwise the user will be
// limited forever.
func (l *Limiter) SetPunishmentDuration(d time.Duration) {
	l.punishment = d
}

// SetMaxMessageCount sets the possible messages count in the
// antifloodwait amout of time (which is `l.timeout`).
// in that period of time, chat (or user) needs to send less than
// this much message, otherwise they will be limited by this limiter
// and so as a result of that their messages will be ignored by the bot.
func (l *Limiter) SetMaxMessageCount(count int) {
	l.maxCount = count
}

// SetMaxCacheDuration will set the max duration
func (l *Limiter) SetMaxCacheDuration(d time.Duration) {
	l.maxTimeout = d
}

// isException will check and see if msg can be ignore because
// it's id is in the exception list or not. This method's usage
// is internal-only.
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

// checker should be run in a new goroutine as it blocks its goroutine
// with a for-loop. This method's duty is to clear the old user's status
// from the cache using `l.maxTimeout` parameter.
func (l *Limiter) checker() {
	for l.isEnabled && !l.isStopped {
		time.Sleep(l.maxTimeout)

		// added this checker just in-case so we can
		// prevent the panics in the future.
		if l.userMap == nil || l.mutex == nil {
			return
		}

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
