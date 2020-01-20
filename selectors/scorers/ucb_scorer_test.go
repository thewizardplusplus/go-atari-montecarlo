package scorers

import (
	"math"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestUCBScorerScoreNode(
	test *testing.T,
) {
	scorer := UCBScorer{Factor: 3}
	score := scorer.ScoreNode(&tree.Node{
		Parent: &tree.Node{
			Children: tree.NodeGroup{
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
		},
		State: tree.NodeState{
			GameCount: 10,
			WinCount:  2,
		},
	})

	roundedScore := math.Floor(score*100) / 100
	if roundedScore != 1.84 {
		test.Fail()
	}
}
