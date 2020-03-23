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
	var waiter sync.WaitGroup
	waiter.Add(len(nodes))

	stateBuffer :=
		make(chan tree.NodeState, len(nodes))
	for _, node := range nodes {
		go func(node *tree.Node) {
			defer waiter.Done()

			stateBuffer <- simulator.Simulator.
				Simulate(node)
		}(node)
	}

	waiter.Wait()
	close(stateBuffer)

	var states []tree.NodeState
	for state := range stateBuffer {
		states = append(states, state)
	}

	return states
}
