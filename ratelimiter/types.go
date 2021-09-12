package ratelimiter

import (
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

// UserStatus is the status of a user in the map.
type UserStatus struct {
	// Last field is the last time that we received a message
	// from the user.
	Last time.Time

	// limited will be true if and only if the current user is
	// banned in the limiter.
	limited bool

	// count is the counts of the messages of the user received
	// by limiter.
	count int
}

// Limiter is the main struct of this library.
type Limiter struct {
	mutex *sync.Mutex
	// IsEnable will be true if and only if the limiter is enabled
	// and should check for the incomming messages.
	isEnabled bool

	// IsStopped will be false when the limiter is stopped.
	isStopped bool

	// userMap is a map of user status with their user id
	// as its key (int64).
	userMap map[int64]*UserStatus

	// trigger function will be runned when a user is limited
	// by the limiter. It should be set by user, users can do everything
	// they want in this function, such as logging the person's id who
	// has been limited by the limiter, etc...
	trigger handlers.Response

	filter filters.Message

	handler handlers.Response

	exceptions   []filters.Message
	exceptionIDs []int64

	timeout    time.Duration
	maxTimeout time.Duration
	maxCount   int

	IgnoreMediaGroup bool
	ConsiderUser     bool
}
