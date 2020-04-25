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
) []tree.NodeState {
	states :=
		make([]tree.NodeState, len(nodes))

	var waiter sync.WaitGroup
	for index, node := range nodes {
		waiter.Add(1)

		go func(index int, node *tree.Node) {
			defer waiter.Done()

			states[index] = simulator.Simulator.
				Simulate(node)
		}(index, node)
	}
	waiter.Wait()

	return states
}
