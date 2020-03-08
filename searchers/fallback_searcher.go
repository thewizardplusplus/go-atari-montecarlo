package searchers

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// FallbackSearcher ...
type FallbackSearcher struct {
	PrimarySearcher  Searcher
	FallbackSearcher Searcher
}

// SearchMove ...
func (searcher FallbackSearcher) SearchMove(
	root *tree.Node,
) (*tree.Node, error) {
	node, err := searcher.PrimarySearcher.
		SearchMove(root)
	if err != nil {
		node, err = searcher.FallbackSearcher.
			SearchMove(root)
	}
	return node, err
}
