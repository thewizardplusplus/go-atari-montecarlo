package tree

// NodeGroup ...
type NodeGroup []*Node

// TotalGameCount ...
func (nodes NodeGroup) TotalGameCount() int {
	var count int
	for _, node := range nodes {
		count += node.State.GameCount
	}

	return count
}
