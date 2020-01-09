package tree

// Node ...
type Node struct {
	Parent *Node
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
