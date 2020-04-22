package tree

import (
	"math"

	models "github.com/thewizardplusplus/go-atari-models"
)

// NodeState ...
type NodeState struct {
	GameCount int
	WinCount  int
}

// NewNodeState ...
//
// Passed error should be
// models.ErrAlreadyLoss or
// models.ErrAlreadyWin only.
//
// Otherwize the function will panic.
func NewNodeState(err error) NodeState {
	state := NodeState{
		GameCount: 1,
	}
	switch err {
	case models.ErrAlreadyLoss:
	case models.ErrAlreadyWin:
		state.WinCount = 1
	default:
		panic(
			"tree.NewNodeState: unsupported error",
		)
	}

	return state
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

// Update ...
func (state *NodeState) Update(
	another NodeState,
) {
	state.GameCount += another.GameCount
	state.WinCount += another.WinCount
}
