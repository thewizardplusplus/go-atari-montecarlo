package scorers

import (
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestWinRateScorerScoreNode(
	test *testing.T,
) {
	var scorer WinRateScorer
	score := scorer.ScoreNode(
		&tree.Node{
			State: tree.NodeState{
				GameCount: 10,
				WinCount:  1,
			},
		},
		[]*tree.Node{
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
		},
	)

	if score != 0.1 {
		test.Fail()
	}
}
