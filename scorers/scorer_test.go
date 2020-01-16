package scorers

import (
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestUCBScorerTotalGameCount(
	test *testing.T,
) {
	var scorer UCBScorer
	count :=
		scorer.totalGameCount([]*tree.Node{
			&tree.Node{
				State: tree.NodeState{
					GameCount: 10,
					WinCount:  1,
				},
			},
			&tree.Node{
				State: tree.NodeState{
					GameCount: 10,
					WinCount:  2,
				},
			},
		})

	if count != 20 {
		test.Fail()
	}
}
