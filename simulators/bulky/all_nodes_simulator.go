package bulky

import (
	"sync"

	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// AllNodesSimulator ...
type AllNodesSimulator struct {
	Simulator simulators.Simulator
}

// Simulate ...
func (simulator AllNodesSimulator) Simulate(
	nodes tree.NodeGroup,
) tree.NodeState {
	var waiter sync.WaitGroup
	waiter.Add(len(nodes))

	states :=
		make(chan tree.NodeState, len(nodes))
	for _, node := range nodes {
		node := node
		go func() {
			defer waiter.Done()

			states <- simulator.Simulator.
				Simulate(node)
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
