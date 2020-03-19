package builders

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// TreeBuilder ...
type TreeBuilder struct {
	NodeSelector tree.NodeSelector
	Simulator    simulators.Simulator
}

// Pass ...
func (builder TreeBuilder) Pass(
	root *tree.Node,
) {
	leaf := root.
		SelectLeaf(builder.NodeSelector).
		ExpandLeaf()[0]
	nextColor := leaf.Move.Color.Negative()
	state := builder.Simulator.Simulate(
		leaf.Board,
		nextColor,
	)
	leaf.UpdateState(state.Invert())
}
