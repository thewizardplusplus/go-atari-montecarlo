package searchers_test

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func search(
	board models.Board,
	color models.Color,
) (move models.Move, ok bool) {
	root :=
		tree.NewNode(board, color.Negative())
	searcher := searchers.MoveSearcher{
		Builder:      nil,
		NodeSelector: nil,
	}
	node, ok := searcher.SearchMove(root)
	if !ok {
		return models.Move{}, false
	}

	return node.Move, true
}
