package tree

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// NodeSelector ...
type NodeSelector interface {
	SelectNode(nodes NodeGroup) *Node
}

// Node ...
type Node struct {
	Parent   *Node
	Move     models.Move
	State    NodeState
	Children NodeGroup
}

// AddResult ...
func (node *Node) AddResult(
	result GameResult,
) {
	node.State.AddResult(result)
	if node.Parent != nil {
		parentResult := result.Invert()
		node.Parent.AddResult(parentResult)
	}
}

// SelectLeaf ...
func (node *Node) SelectLeaf(
	selector NodeSelector,
) *Node {
	for len(node.Children) > 0 {
		node = selector.SelectNode(node.Children)
	}

	return node
}
