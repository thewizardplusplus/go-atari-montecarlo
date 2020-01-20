package tree

// NodeState ...
type NodeState struct {
	GameCount int
	WinCount  int
}

// WinRate ...
func (state NodeState) WinRate() float64 {
	return float64(state.WinCount) /
		float64(state.GameCount)
}

// AddResult ...
func (state *NodeState) AddResult(
	result GameResult,
) {
	state.GameCount++
	if result == Win {
		state.WinCount++
	}
}
