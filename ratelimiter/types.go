// ratelimiter Project
// Copyright (C) 2021 ALiwoto and other Contributors
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ratelimiter

import (
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
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

	custom *customIgnore
}

type customIgnore struct {
	startTime       time.Time
	duration        time.Duration
	ignoreException bool
}

// Limiter is the main struct of this library.
type Limiter struct {
	mutex *sync.Mutex
	// IsEnable will be true if and only if the limiter is enabled
	// and should check for the incoming messages.
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
	triggers []handlers.Response

	filter filters.Message

	handler handlers.Response

	// msgHandler is the original message handler of this limiter.
	// it should remain private.
	msgHandler *handlers.Message

	allHandlers []ext.Handler

	exceptions        []filters.Message
	conditions        []filters.Message
	exceptionIDs      []int64
	ignoredExceptions []int64

	// timeout is the floodwait checking time. a user is allowed to
	// send `maxCount` messages per `timeout`.
	timeout time.Duration

	// maxTimeout is the maximum time out of clearing user status
	// cache in the memory.
	maxTimeout time.Duration

	// punishment is the necessary time a user needs to spend after
	// being limiter as its punishment; the user will be freed after
	// this time is passed.
	punishment time.Duration

	// maxCount is the maximum number of messages we can accept from the
	// user in `timeout` amount of time; if the user sends more than
	// this much message, it will be limited and so the bot will ignore
	// their messages.
	maxCount int

	// IgnoreMediaGroup should be set to true when we have to ignore
	// album messages (such as album musics, album photos, etc...) and
	// don't check them at all.
	// default value for this field is true.
	IgnoreMediaGroup bool

	// TextOnly should be set to true when we have to ignore
	// media messages (such as photos, videos, audios, etc...) and
	// don't check them at all.
	// If your bot has nothing to do with media messages, you can set
	// this to true.
	TextOnly bool

	// IsStrict will tell the limiter whether it should act more strict
	// or not. If this value is set to `true`, the user should NOT send
	// any messages to the bot until it's limit time is completely over.
	// otherwise the limitation will remain on the user until it stops
	// sending any messages to the bot.
	// (A truly bad way of handling antifloodwait... we recommend not to
	// set this value to `true`, unless it's very very necessary).
	IsStrict bool

	// ConsiderUser will be true when the limiter needs to consider users
	// for their checking the messages. so the user's ID will be used as key
	// to access the map.
	ConsiderUser bool

	// ConsiderInline fields will determine whether we need to
	ConsiderInline bool
}

// LimiterConfig is the config type of the limiter.
type LimiterConfig struct {
	ConsiderChannel  bool
	ConsiderUser     bool
	ConsiderEdits    bool
	IgnoreMediaGroup bool
	TextOnly         bool
	IsStrict         bool
	HandlerGroups    []int
	ConsiderInline   bool
}
