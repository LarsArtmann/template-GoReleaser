package jobs

import (
	"time"
)

const (
	defaultEstimatedTime   = 5 * time.Second
	shortEstimatedTime     = 3 * time.Second
	veryShortEstimatedTime = 2 * time.Second
	defaultMaxRetries      = 2
	defaultWorkflowTimeout = 30 * time.Minute
	defaultJobTimeout      = 5 * time.Minute
)
