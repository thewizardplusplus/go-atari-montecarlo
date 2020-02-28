package searchers

import (
	"errors"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// ...
var (
	ErrNotFoundPreviousMove = errors.New(
		"not found previous move",
	)
)

// Searcher ...
type Searcher interface {
	SearchMove(
		root *tree.Node,
	) (*tree.Node, error)
}

// ReusedSearcher ...
type ReusedSearcher struct {
	searcher     Searcher
	previousMove *tree.Node
}

// NewReusedSearcher ...
func NewReusedSearcher(
	searcher Searcher,
) ReusedSearcher {
	return ReusedSearcher{searcher: searcher}
}

// SearchMove ...
func (searcher ReusedSearcher) SearchMove(
	root *tree.Node,
) (*tree.Node, error) {
	if searcher.previousMove != nil {
		node, ok :=
			searcher.searchPreviousMove(root)
		if !ok {
			return nil, ErrNotFoundPreviousMove
		}

		root = node
	}

	node, err :=
		searcher.searcher.SearchMove(root)
	if err != nil {
		return nil, err
	}

	searcher.previousMove = node

	return node, nil
}

func (
	searcher ReusedSearcher,
) searchPreviousMove(
	sample *tree.Node,
) (node *tree.Node, ok bool) {
	nodes := searcher.previousMove.Children
	for _, node := range nodes {
		if node.Move == sample.Move {
			return node, true
		}
	}

	return nil, false
}
