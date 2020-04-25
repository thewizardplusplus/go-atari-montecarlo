package bulky

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/parallel"
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
	packedStates := parallel.Run(
		len(nodes),
		func(index int) (result interface{}) {
			return simulator.Simulator.
				Simulate(nodes[index])
		},
	)

	states :=
		make([]tree.NodeState, 0, len(nodes))
	for _, packedState := range packedStates {
		states = append(
			states,
			packedState.(tree.NodeState),
		)
	}

	return states
}
