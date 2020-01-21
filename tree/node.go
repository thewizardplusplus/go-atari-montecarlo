package tree

import (
	"errors"

	models "github.com/thewizardplusplus/go-atari-models"
)

// ...
var (
	ErrNoMoves = errors.New("no moves")
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

// ExpandLeaf ...
func (
	node *Node,
) ExpandLeaf() (*Node, error) {
	if node.State.GameCount == 0 {
		return node, nil
	}

	nextColor := node.Move.Color.Negative()
	moves := node.Board.Moves(nextColor)
	if len(moves) == 0 {
		return nil, ErrNoMoves
	}

	var children NodeGroup
	for _, move := range moves {
		nextBoard := node.Board.ApplyMove(move)
		child := &Node{
			Parent: node,
			Move:   move,
			Board:  nextBoard,
		}
		children = append(children, child)
	}

	return children[0], nil
}
