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

// ShallowCopy ...
//
// It copies only the move and the board.
func (node *Node) ShallowCopy() *Node {
	return &Node{
		Move:  node.Move,
		Board: node.Board,
	}
}

// UpdateState ...
func (node *Node) UpdateState(
	state NodeState,
) {
	node.State.Update(state)
	if node.Parent != nil {
		parentState := state.Invert()
		node.Parent.UpdateState(parentState)
	}
}

// MergeChildren ...
//
// It merges only states of children.
//
// If the argument hasn't children,
// then this method does nothing.
//
// If argument children don't contain
// any node, then the latter isn't updated.
//
// If argument children contain
// any additional node, then the latter
// is ignored.
//
// If this node hasn't children, it borrows
// argument children.
func (node *Node) MergeChildren(
	another *Node,
) {
	if len(node.Children) != 0 {
		node.Children.Merge(another.Children)
		return
	}

	// node hasn't children;
	// borrow argument children
	node.Children = another.Children
	for _, child := range node.Children {
		child.Parent = node
	}
}

// SelectLeaf ...
func (node *Node) SelectLeaf(
	selector NodeSelector,
) *Node {
	for len(node.Children) > 0 {
		node =
			selector.SelectNode(node.Children)
	}

	return node
}

// ExpandLeaf ...
func (node *Node) ExpandLeaf() NodeGroup {
	if node.State.GameCount == 0 {
		return NodeGroup{node}
	}

	moves, err :=
		node.Board.LegalMoves(node.Move)
	if err != nil {
		// no moves or an already finished game
		return NodeGroup{node}
	}

	node.Children = NewNodeGroup(node, moves)
	return node.Children
}
