package searchers

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// MoveSearcher ...
type MoveSearcher struct {
	Builder      builders.Builder
	NodeSelector tree.NodeSelector
}

// SearchMove ...
func (searcher MoveSearcher) SearchMove(
	root *tree.Node,
) (node *tree.Node, ok bool) {
	searcher.Builder.Pass(root)
	if len(root.Children) == 0 {
		return nil, false
	}

	node = searcher.NodeSelector.
		SelectNode(root.Children)
	return node, true
}
