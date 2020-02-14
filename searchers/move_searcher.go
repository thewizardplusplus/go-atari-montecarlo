package searchers

import (
	"log"

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
		log.Print("no children")
		log.Println("")

		return nil, false
	}

	for i, c := range root.Children {
		log.Printf("child #%d", i)
		log.Printf("move %v", c.Move)
		log.Printf("state %v", c.State)
	}
	log.Println("")

	node = searcher.NodeSelector.
		SelectNode(root.Children)
	return node, true
}
