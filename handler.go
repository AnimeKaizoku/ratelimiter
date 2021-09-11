package ratelimiter

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *RateLimiter) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveMessage == nil || ctx.EffectiveMessage.From == nil {
		return ext.ContinueGroups
	}

	key := c.KeyGenerator(ctx)

	c.mutex.Lock()
	
	c.hits[key] = c.hits[key] + 1
	if c.hits[key] > c.Limit {
		c.mutex.Unlock()
		return ext.EndGroups
	}

	c.mutex.Unlock()
	return ext.ContinueGroups
}

func (c *RateLimiter) CheckUpdate(b *gotgbot.Bot, u *gotgbot.Update) bool {
	return u.CallbackQuery != nil || u.Message != nil
}

func (c *RateLimiter) Name() string {
	return fmt.Sprintf("ratelimiter_%p", c.HandleUpdate)
}
