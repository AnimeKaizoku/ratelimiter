<!--
 *
 * This file is part of ratelimiter-gotgbot (https://github.com/gotgbot/ratelimiter).
 * Copyright (c) 2021 ALiwoto, Contributors.
 *
 * This library is free: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License License as published by
 * the Free Software Foundation, version 3.
 *
 * This library is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
-->


# <p align="center"> Rate Limiter - gotgbot </p>

<p align="center">
	<a href="https://github.com/gotgbot/ratelimiter">
		<img src="./logo.png" alt="ratelimiter-Logo">
	</a>
</p>

> Name:		Rate Limiter			\
> Version:	v1.0.7					\
> Edit:		2 Feb 2021				\
> By:		ALiwoto and Contributors (C)	

[![Go Reference](https://pkg.go.dev/badge/github.com/gotgbot/ratelimiter.svg)](https://pkg.go.dev/github.com/gotgbot/ratelimiter) [![Go-linux](https://github.com/gotgbot/ratelimiter/actions/workflows/go-linux.yml/badge.svg)](https://github.com/gotgbot/ratelimiter/actions/workflows/go-linux.yml) [![Go-macos](https://github.com/gotgbot/ratelimiter/actions/workflows/go-macos.yml/badge.svg)](https://github.com/gotgbot/ratelimiter/actions/workflows/go-macos.yml) [![Go-windows](https://github.com/gotgbot/ratelimiter/actions/workflows/go-windows.yml/badge.svg)](https://github.com/gotgbot/ratelimiter/actions/workflows/go-windows.yml)

<hr/>

## How to use

```go
import "github.com/gotgbot/ratelimiter/ratelimiter"


func loadLimiter(d *ext.Dispatcher) {
	limiter = ratelimiter.NewLimiter(d, &ratelimiter.LimiterConfig{
		ConsiderChannel:  false,
		ConsiderUser:     true,
		ConsiderEdits:    false,
		IgnoreMediaGroup: true,
		TextOnly:         false,
		HandlerGroups:    []int{0, 1, 2},
	})

	// 14 messages per 6 seconds
	limiter.SetFloodWaitTime(6 * time.Second)
	limiter.SetMaxMessageCount(14)

	// add sudo users as exceptions, so they don't get rate-limited by library
	limiter.AddExceptionID(sudoUsers...)

	limiter.Start()
}
```

<hr/>

## Helpful links:

- [Support group on telegram](https://t.me/KaizokuBots)
- [Contact maintainer on telegram](https://t.me/Falling_inside_The_Black)

