// ratelimiter Project
// Copyright (C) 2021~2022 ALiwoto and other Contributors
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ratelimiter

import "time"

const (
	DefaultTimeout        = 4 * time.Second
	DefaultPunishmentTime = 4 * time.Minute
	DefaultMaxTimeout     = 30 * time.Minute
	DefaultMessageCount   = 15
)
