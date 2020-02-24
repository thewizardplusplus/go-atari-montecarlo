package terminators

// BuildingTerminator ...
type BuildingTerminator interface {
	IsSearchTerminated(pass int) bool
}

// GroupTerminator ...
type GroupTerminator struct {
	terminators []BuildingTerminator
}

// NewGroupTerminator ...
func NewGroupTerminator(
	terminators ...BuildingTerminator,
) GroupTerminator {
	return GroupTerminator{terminators}
}

// IsSearchTerminated ...
func (
	group GroupTerminator,
) IsSearchTerminated(pass int) bool {
	terminators := group.terminators
	for _, terminator := range terminators {
		if terminator.IsSearchTerminated(pass) {
			return true
		}
	}

	return false
}
