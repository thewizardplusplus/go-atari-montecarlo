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
	got := selector.SelectNode(tree.NodeGroup{
		&tree.Node{
			State: tree.NodeState{
				GameCount: 2,
				WinCount:  1,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 4,
				WinCount:  3,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 6,
				WinCount:  5,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 8,
				WinCount:  9,
			},
		},
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  9,
			},
		},
	})

	want := &tree.Node{
		State: tree.NodeState{
			GameCount: 4,
			WinCount:  3,
		},
	}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}
