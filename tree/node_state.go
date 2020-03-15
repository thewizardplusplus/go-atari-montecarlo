package tree

import (
	"math"
)

// NodeState ...
type NodeState struct {
	GameCount int
	WinCount  int
}

// WinRate ...
func (state NodeState) WinRate() float64 {
	if state.GameCount == 0 {
		return math.Inf(+1)
	}

	return float64(state.WinCount) /
		float64(state.GameCount)
}

// Invert ...
func (state NodeState) Invert() NodeState {
	return NodeState{
		GameCount: state.GameCount,
		WinCount: state.GameCount -
			state.WinCount,
	}
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

// Update ...
func (state *NodeState) Update(
	another NodeState,
) {
	state.GameCount += another.GameCount
	state.WinCount += another.WinCount
}
