package searchers

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockSearcher struct {
	searchMove func(
		root *tree.Node,
	) (*tree.Node, error)
}

func (searcher MockSearcher) SearchMove(
	root *tree.Node,
) (*tree.Node, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(root)
}

func TestNewReusedSearcher(
	test *testing.T,
) {
	var innerSearcher MockSearcher
	searcher :=
		NewReusedSearcher(innerSearcher)

	if !reflect.DeepEqual(
		searcher.searcher,
		innerSearcher,
	) {
		test.Fail()
	}
	if searcher.previousMove != nil {
		test.Fail()
	}
}
