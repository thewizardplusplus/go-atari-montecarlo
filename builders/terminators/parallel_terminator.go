package terminators

import (
	"sync"
)

// ParallelTerminator ...
type ParallelTerminator struct {
	locker     sync.RWMutex
	terminator BuildingTerminator
}

// NewParallelTerminator ...
func NewParallelTerminator(
	terminator BuildingTerminator,
) *ParallelTerminator {
	return &ParallelTerminator{
		terminator: terminator,
	}
}

// IsBuildingTerminated ...
func (
	terminator *ParallelTerminator,
) IsBuildingTerminated(pass int) bool {
	terminator.locker.RLock()
	defer terminator.locker.RUnlock()

	return terminator.terminator.
		IsBuildingTerminated(pass)
}

// Reset ...
func (
	terminator *ParallelTerminator,
) Reset() {
	terminator.locker.Lock()
	defer terminator.locker.Unlock()

	terminator.terminator.Reset()
}
