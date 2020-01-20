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
	Factor float64
}

// ScoreNode ...
func (scorer UCBScorer) ScoreNode(
	node *tree.Node,
) float64 {
	x := node.State.WinRate()
	n := scorer.
		totalGameCount(node.Parent.Children)
	shift := scorer.Factor * math.Sqrt(
		math.Log(float64(n))/
			float64(node.State.GameCount),
	)
	return x + shift
}

func (scorer UCBScorer) totalGameCount(
	nodes tree.NodeGroup,
) int {
	var count int
	for _, node := range nodes {
		count += node.State.GameCount
	}

	return count
}
