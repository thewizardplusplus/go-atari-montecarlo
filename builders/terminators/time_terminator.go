package terminators

import (
	"time"
)

// Clock ...
type Clock func() time.Time

// TimeTerminator ...
type TimeTerminator struct {
	clock           Clock
	maximalDuration time.Duration
	startTime       time.Time
}

// NewTimeTerminator ...
func NewTimeTerminator(
	clock Clock,
	maximalDuration time.Duration,
) *TimeTerminator {
	return &TimeTerminator{
		clock:           clock,
		maximalDuration: maximalDuration,
		startTime:       clock(),
	}
}

// IsBuildingTerminated ...
func (terminator TimeTerminator) IsBuildingTerminated(pass int) bool {
	currentTime := terminator.clock()
	elapsedTime := currentTime.Sub(terminator.startTime)
	return elapsedTime >= terminator.maximalDuration
}
