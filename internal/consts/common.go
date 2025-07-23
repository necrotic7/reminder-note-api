package consts

import "time"

// Default
const (
	PageSize      = 10
	Timeout       = 5 * time.Second
	Semaphore     = 10
	RetryInterval = 3 * time.Second
)
