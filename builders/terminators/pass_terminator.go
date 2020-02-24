package terminators

// PassTerminator ...
type PassTerminator struct {
	maximalPass int
}

// NewPassTerminator ...
func NewPassTerminator(
	maximalPass int,
) PassTerminator {
	return PassTerminator{maximalPass}
}

// IsSearchTerminated ...
func (
	terminator PassTerminator,
) IsSearchTerminated(pass int) bool {
	return pass >= terminator.maximalPass
}
