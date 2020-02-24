package terminators

// BuildingTerminator ...
type BuildingTerminator interface {
	IsBuildingTerminated(pass int) bool
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

// IsBuildingTerminated ...
func (
	group GroupTerminator,
) IsBuildingTerminated(pass int) bool {
	terminators := group.terminators
	for _, terminator := range terminators {
		if terminator.
			IsBuildingTerminated(pass) {
			return true
		}
	}

	return false
}
