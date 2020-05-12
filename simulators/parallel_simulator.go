package simulators

import (
	syncutils "github.com/thewizardplusplus/go-atari-montecarlo/sync-utils"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// Simulator ...
type Simulator interface {
	Simulate(root *tree.Node) tree.NodeState
}

// ParallelSimulator ...
type ParallelSimulator struct {
	Simulator   Simulator
	Concurrency int
}

// Simulate ...
func (simulator ParallelSimulator) Simulate(
	root *tree.Node,
) tree.NodeState {
	states := syncutils.ParallelRun(
		simulator.Concurrency,
		func(index int) (result interface{}) {
			return simulator.Simulator.
				Simulate(root)
		},
	)

	var generalState tree.NodeState
	for _, state := range states {
		generalState.
			Update(state.(tree.NodeState))
	}

	return generalState
}
