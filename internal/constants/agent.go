package constants

import "time"

const (
	MaxDelayForWaitingServer      = 10 * time.Second
	IncreaseDelayForWaitingServer = 2 * time.Second
)
