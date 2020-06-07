package scorers

import (
	"math"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// UCBScorer ...
//
// It implements the Upper Confidence Bound algorithm.
//
type UCBScorer struct {
	Factor float64
}

// ScoreNode ...
func (scorer UCBScorer) ScoreNode(node *tree.Node) float64 {
	x := node.State.WinRate()
	if x == math.Inf(+1) {
		return x
	}

	shift := scorer.Factor *
		math.Sqrt(math.Log(gameCount(node.Parent))/gameCount(node))
	return x + shift
}

func gameCount(node *tree.Node) float64 {
	return float64(node.State.GameCount)
}
