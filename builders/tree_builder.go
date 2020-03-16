package builders

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// Simulator ...
type Simulator interface {
	Simulate(
		board models.Board,
		color models.Color,
	) tree.NodeState
}

// TreeBuilder ...
type TreeBuilder struct {
	NodeSelector tree.NodeSelector
	Simulator    Simulator
}

// Pass ...
func (builder TreeBuilder) Pass(
	root *tree.Node,
) {
	leaf := root.
		SelectLeaf(builder.NodeSelector).
		ExpandLeaf()
	nextColor := leaf.Move.Color.Negative()
	state := builder.Simulator.Simulate(
		leaf.Board,
		nextColor,
	)
	leaf.UpdateState(state.Invert())
}
