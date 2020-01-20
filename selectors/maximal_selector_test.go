package selectors

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockNodeScorer struct{}

func (scorer MockNodeScorer) ScoreNode(
	node *tree.Node,
) float64 {
	return node.State.WinRate()
}

func TestMaximalSelectorSelectNode(
	test *testing.T,
) {
	selector := MaximalSelector{
		NodeScorer: MockNodeScorer{},
	}
	got := selector.SelectNode([]*tree.Node{
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  1,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  3,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  5,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  4,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  2,
			},
		},
	})

	want := &tree.Node{
		State: tree.NodeState{
			GameCount: 10,
			WinCount:  5,
		},
	}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}
