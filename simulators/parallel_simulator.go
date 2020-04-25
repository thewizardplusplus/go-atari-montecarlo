package simulators

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/parallel"
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
	states := parallel.Run(
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
