package tree

// NodeState ...
type NodeState struct {
	GameCount int
	WinCount  int
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
