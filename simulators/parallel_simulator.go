package simulators

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// Simulator ...
type Simulator interface {
	Simulate(
		board models.Board,
		color models.Color,
	) tree.NodeState
}

// ParallelSimulator ...
type ParallelSimulator struct {
	Simulator   Simulator
	Concurrency int
}

// Simulate ...
func (simulator ParallelSimulator) Simulate(
	board models.Board,
	color models.Color,
) tree.NodeState {
	return tree.NodeState{}
}
