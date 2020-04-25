package simulators

import (
	"sync"

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
	states := make(
		[]tree.NodeState,
		simulator.Concurrency,
	)

	var waiter sync.WaitGroup
	for i := 0; i < len(states); i++ {
		waiter.Add(1)

		go func(i int) {
			defer waiter.Done()

			states[i] = simulator.Simulator.
				Simulate(root)
		}(i)
	}
	waiter.Wait()

	var generalState tree.NodeState
	for _, state := range states {
		generalState.Update(state)
	}

	return generalState
}
