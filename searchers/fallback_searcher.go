package searchers

import (
	models "github.com/thewizardplusplus/go-atari-models"
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
	switch err {
	case nil:
		return node, nil
	case models.ErrAlreadyWin,
		models.ErrAlreadyLoss:
		return nil, err
	}

	return searcher.FallbackSearcher.
		SearchMove(root)
}
