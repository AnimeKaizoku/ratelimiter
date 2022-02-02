// ratelimiter Project
// Copyright (C) 2021~2022 ALiwoto and other Contributors
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ratelimiter

import "time"

const (
	DEFAULT_TIME        = 4 * time.Second
	DEFAULT_PUNISHMENT  = 4 * time.Minute
	DEFAULT_MAX_TIMEOUT = 30 * time.Minute
	DEFAULT_COUNT       = 15
)
