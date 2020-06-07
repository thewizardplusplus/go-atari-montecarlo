package terminators

// PassTerminator ...
type PassTerminator struct {
	maximalPass int
}

// NewPassTerminator ...
func NewPassTerminator(maximalPass int) PassTerminator {
	return PassTerminator{maximalPass}
}

// IsBuildingTerminated ...
func (terminator PassTerminator) IsBuildingTerminated(pass int) bool {
	return pass >= terminator.maximalPass
}
