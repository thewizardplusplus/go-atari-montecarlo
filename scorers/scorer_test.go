package scorers

import (
	"math"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestUCBScorerScoreNode(
	test *testing.T,
) {
	scorer := UCBScorer{Factor: 2}
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

	roundedScore := math.Floor(score*100) / 100
	if roundedScore != 1.19 {
		test.Log(score)
		test.Fail()
	}
}

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
