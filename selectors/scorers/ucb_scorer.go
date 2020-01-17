package scorers

import (
	"math"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// UCBScorer ...
//
// It implements
// the Upper Confidence Bound algorithm.
type UCBScorer struct {
	WinRateScorer

	Factor float64
}

// ScoreNode ...
func (scorer UCBScorer) ScoreNode(
	node *tree.Node,
	siblings []*tree.Node,
) float64 {
	x := scorer.WinRateScorer.
		ScoreNode(node, siblings)
	n := scorer.totalGameCount(siblings)
	shift := scorer.Factor * math.Sqrt(
		math.Log(float64(n))/
			float64(node.State.GameCount),
	)
	return x + shift
}

func (scorer UCBScorer) totalGameCount(
	siblings []*tree.Node,
) int {
	var count int
	for _, sibling := range siblings {
		count += sibling.State.GameCount
	}

	return count
}
