package ratelimiter

import (
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// The function which generates a unique ID for the given context which is later used to count hits.
type KeyGenerator func(ctx *ext.Context) string

type RateLimiter struct {
	// `Limit` is the number of messages that are allowed to be received in the `TimeFrame`
	Limit     int
	TimeFrame time.Duration
	KeyGenerator
	hits        map[string]int
	mutex       *sync.Mutex
	initialized bool
}

// Initializes a new rate limiter which can be later added as a handler.
func NewRatelimiter(timeFrame time.Duration, limit int, keyGenerator KeyGenerator) *RateLimiter {
	r := &RateLimiter{TimeFrame: timeFrame, Limit: limit, KeyGenerator: keyGenerator}

	r.Initialize()
	return r
}
