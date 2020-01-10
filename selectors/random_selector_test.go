package selectors

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestRandomSelectorSelectNode(
	test *testing.T,
) {
	// make the random generator deterministic
	// for test reproducibility
	rand.Seed(1)

	var selector RandomSelector
	got := selector.SelectNode([]*tree.Node{
		&tree.Node{
			State: tree.NodeState{
				GameCount: 1,
				WinCount:  2,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 3,
				WinCount:  4,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 5,
				WinCount:  6,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 7,
				WinCount:  8,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 9,
				WinCount:  10,
			},
		},
	})

	want := &tree.Node{
		State: tree.NodeState{
			GameCount: 3,
			WinCount:  4,
		},
	}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}
