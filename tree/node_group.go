package tree

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// NodeGroupConfiguration ...
type NodeGroupConfiguration struct {
	parent   *Node
	useBoard bool
	board    models.Board
}

// NodeGroupOption ...
type NodeGroupOption func(
	configuration *NodeGroupConfiguration,
)

// WithParent ...
func WithParent(
	parent *Node,
) NodeGroupOption {
	return func(
		configuration *NodeGroupConfiguration,
	) {
		configuration.parent = parent
	}
}

// WithBoard ...
func WithBoard(
	board models.Board,
) NodeGroupOption {
	return func(
		configuration *NodeGroupConfiguration,
	) {
		configuration.useBoard = true
		configuration.board = board
	}
}

// NodeGroup ...
type NodeGroup []*Node

// NewNodeGroup ...
func NewNodeGroup(
	moves []models.Move,
	options ...NodeGroupOption,
) NodeGroup {
	var configuration NodeGroupConfiguration
	for _, option := range options {
		option(&configuration)
	}

	var nodes NodeGroup
	for _, move := range moves {
		node := &Node{
			Parent: configuration.parent,
			Move:   move,
		}

		if configuration.useBoard {
			nextBoard :=
				configuration.board.ApplyMove(move)
			node.Board = nextBoard
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// TotalGameCount ...
func (
	nodes NodeGroup,
) TotalGameCount() int {
	var count int
	for _, node := range nodes {
		count += node.State.GameCount
	}

	return count
}

// Merge ...
//
// It merges only states of nodes.
//
// If the argument is nil,
// then this method does nothing.
//
// If the argument doesn't contain any move,
// then the latter isn't updated.
//
// If the argument contains
// any additional move, then the latter
// is ignored.
func (nodes NodeGroup) Merge(
	another NodeGroup,
) {
	anotherStates :=
		make(map[models.Move]NodeState)
	for _, node := range another {
		anotherStates[node.Move] = node.State
	}

	for _, node := range nodes {
		anotherState := anotherStates[node.Move]
		node.UpdateState(anotherState)
	}
}
