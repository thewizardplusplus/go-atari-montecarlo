package searchers

import (
	"errors"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// ...
var (
	ErrFailedBuilding = errors.New(
		"failed building",
	)
)

// MoveSearcher ...
type MoveSearcher struct {
	MoveGenerator models.Generator
	Builder       builders.Builder
	NodeSelector  tree.NodeSelector
}

// SearchMove ...
//
// Returned error can be
// models.ErrAlreadyLoss,
// models.ErrAlreadyWin or
// ErrFailedBuilding only.
func (searcher MoveSearcher) SearchMove(
	root *tree.Node,
) (*tree.Node, error) {
	_, err := searcher.MoveGenerator.
		LegalMoves(root.Board, root.Move)
	if err != nil {
		return nil, err
	}

	searcher.Builder.Pass(root)
	if len(root.Children) == 0 {
		return nil, ErrFailedBuilding
	}

	node := searcher.NodeSelector.
		SelectNode(root.Children)
	return node, nil
}
