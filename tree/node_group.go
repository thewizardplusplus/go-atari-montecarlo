package tree

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// NodeGroup ...
type NodeGroup []*Node

// NewNodeGroup ...
func NewNodeGroup(parent *Node, moves []models.Move) NodeGroup {
	var nodes NodeGroup
	for _, move := range moves {
		nextStorage := parent.Storage.ApplyMove(move)
		node := &Node{
			Parent:  parent,
			Move:    move,
			Storage: nextStorage,
		}
		nodes = append(nodes, node)
	}

	return nodes
}

// Merge ...
//
// It merges only states of nodes.
//
// If the argument is nil, then this method does nothing.
//
// If the argument doesn't contain any move, then the latter isn't updated.
//
// If the argument contains any additional move, then the latter is ignored.
//
func (nodes NodeGroup) Merge(another NodeGroup) {
	anotherStates := make(map[models.Move]NodeState)
	for _, node := range another {
		anotherStates[node.Move] = node.State
	}

	for _, node := range nodes {
		anotherState := anotherStates[node.Move]
		node.UpdateState(anotherState)
	}
}
