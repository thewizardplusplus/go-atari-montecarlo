package selectors

import (
	"math"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// NodeScorer ...
type NodeScorer interface {
	ScoreNode(node *tree.Node) float64
}

// MaximalNodeSelector ...
type MaximalNodeSelector struct {
	NodeScorer NodeScorer
}

// SelectNode ...
func (
	selector MaximalNodeSelector,
) SelectNode(
	nodes tree.NodeGroup,
) *tree.Node {
	var maximum *tree.Node
	maximumScore := math.Inf(-1)
	for _, node := range nodes {
		nodeScore :=
			selector.NodeScorer.ScoreNode(node)
		if nodeScore > maximumScore {
			maximum = node
			maximumScore = nodeScore
		}
	}

	return maximum
}