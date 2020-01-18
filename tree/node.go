package tree

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// Node ...
type Node struct {
	Parent *Node
	Move   models.Move
	State  NodeState
}

// AddResult ...
func (node *Node) AddResult(
	result GameResult,
) {
	node.State.AddResult(result)
	if node.Parent != nil {
		node.Parent.AddResult(result)
	}
}
