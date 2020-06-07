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
func NewGroupTerminator(terminators ...BuildingTerminator) GroupTerminator {
	return GroupTerminator{terminators}
}

// IsBuildingTerminated ...
func (group GroupTerminator) IsBuildingTerminated(pass int) bool {
	for _, terminator := range group.terminators {
		if terminator.IsBuildingTerminated(pass) {
			return true
		}
	}

	return false
}
