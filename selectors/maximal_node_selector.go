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
func (selector MaximalNodeSelector) SelectNode(
	nodes tree.NodeGroup,
) *tree.Node {
	var maximalNode *tree.Node
	maximalNodeScore := math.Inf(-1)
	for _, node := range nodes {
		nodeScore := selector.NodeScorer.ScoreNode(node)
		if nodeScore > maximalNodeScore {
			maximalNode = node
			maximalNodeScore = nodeScore
		}
	}

	return maximalNode
}
