package ratelimiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Initializes the rate limiter. Not required to call if used `NewRateLimiter`.
func (r *RateLimiter) Initialize() {
	if r.initialized {
		return
	}

	r.mutex = &sync.Mutex{}

	if r.Limit == 0 {
		r.Limit = 1
	}

	if r.TimeFrame == 0 {
		r.TimeFrame = time.Second
	}

	if r.KeyGenerator == nil {
		r.KeyGenerator = func(ctx *ext.Context) string {
			return string(strconv.AppendInt([]byte(""), ctx.EffectiveMessage.From.Id, 10))
		}
	}

	go func() {
		for {
			time.Sleep(r.TimeFrame)
			r.mutex.Lock()
			r.hits = make(map[string]int)
			r.mutex.Unlock()
		}
	}()

	r.initialized = true
}
