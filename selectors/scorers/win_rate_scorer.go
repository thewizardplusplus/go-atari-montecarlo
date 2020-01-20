package scorers

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// WinRateScorer ...
type WinRateScorer struct{}

// ScoreNode ...
func (scorer WinRateScorer) ScoreNode(
	node *tree.Node,
	siblings []*tree.Node,
) float64 {
	return node.State.WinRate()
}
