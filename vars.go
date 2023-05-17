package ratelimiter

var (
	DefaultConfig *LimiterConfig = &LimiterConfig{
		ConsiderChannel:  false,
		ConsiderUser:     true,
		ConsiderEdits:    false,
		IgnoreMediaGroup: true,
		TextOnly:         false,
		ConsiderInline:   true,
		IsStrict:         false,
		Timeout:          DefaultTimeout,
		PunishmentTime:   DefaultPunishmentTime,
		MaxTimeout:       DefaultMaxTimeout,
		MessageCount:     DefaultMessageCount,
	}
)
