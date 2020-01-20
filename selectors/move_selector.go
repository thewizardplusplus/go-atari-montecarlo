package selectors

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// NodeSelector ...
type NodeSelector interface {
	SelectNode(nodes tree.NodeGroup) *tree.Node
}

// MoveSelector ...
type MoveSelector struct {
	NodeSelector NodeSelector
}

// SelectMove ...
func (selector MoveSelector) SelectMove(
	moves []models.Move,
) models.Move {
	var nodes []*tree.Node
	for _, move := range moves {
		node := &tree.Node{Move: move}
		nodes = append(nodes, node)
	}

	node :=
		selector.NodeSelector.SelectNode(nodes)
	return node.Move
}
