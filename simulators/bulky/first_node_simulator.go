package bulky

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// FirstNodeSimulator ...
type FirstNodeSimulator struct {
	Simulator simulators.Simulator
}

// Simulate ...
func (simulator FirstNodeSimulator) Simulate(
	nodes tree.NodeGroup,
) []tree.NodeState {
	state := simulator.Simulator.
		Simulate(nodes[0])
	return []tree.NodeState{state}
}
