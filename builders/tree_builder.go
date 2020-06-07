package builders

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// BulkySimulator ...
type BulkySimulator interface {
	// States should correspond to nodes.
	Simulate(nodes tree.NodeGroup) []tree.NodeState
}

// TreeBuilder ...
type TreeBuilder struct {
	NodeSelector  tree.NodeSelector
	MoveGenerator models.Generator
	Simulator     BulkySimulator
}

// Pass ...
func (builder TreeBuilder) Pass(root *tree.Node) {
	leaves := root.
		SelectLeaf(builder.NodeSelector).
		ExpandLeaf(builder.MoveGenerator)
	states := builder.Simulator.Simulate(leaves)
	for index, state := range states {
		leaves[index].UpdateState(state.Invert())
	}
}
