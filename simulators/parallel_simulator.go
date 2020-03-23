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
	var waiter sync.WaitGroup
	waiter.Add(simulator.Concurrency)

	states := make(
		chan tree.NodeState,
		simulator.Concurrency,
	)
	concurrency := simulator.Concurrency
	for i := 0; i < concurrency; i++ {
		go func() {
			defer waiter.Done()

			states <- simulator.Simulator.
				Simulate(root)
		}()
	}

	waiter.Wait()
	close(states)

	var generalState tree.NodeState
	for state := range states {
		generalState.Update(state)
	}

	return generalState
}
