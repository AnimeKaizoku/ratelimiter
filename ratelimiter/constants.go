package ratelimiter

import "time"

const (
	DEFAULT_TIME        = 4 * time.Second
	DEFAULT_MAX_TIMEOUT = 30 * time.Minute
	DEFAULT_COUNT       = 20
)
