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
	Board    models.Board
	State    NodeState
	Children NodeGroup
}

// NewPreliminaryNode ...
//
// It creates the node
// from only the passed board and
// the negated passed color.
func NewPreliminaryNode(
	board models.Board,
	color models.Color,
) *Node {
	return &Node{
		Move:  models.NewPreliminaryMove(color),
		Board: board,
	}
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

// Update ...
func (node *Node) Update(state NodeState) {
	node.State.Update(state)
	if node.Parent != nil {
		parentState := state.Invert()
		node.Parent.Update(parentState)
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

// ExpandLeaf ...
func (node *Node) ExpandLeaf() *Node {
	if node.State.GameCount == 0 {
		return node
	}

	nextColor := node.Move.Color.Negative()
	moves, err := node.Board.
		LegalMoves(nextColor)
	if err != nil {
		// no moves or an already finished game
		return node
	}

	node.Children = NewNodeGroup(
		moves,
		WithParent(node),
		WithBoard(node.Board),
	)

	return node.Children[0]
}
