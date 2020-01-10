package selectors

import (
	"math"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// NodeScorer ...
type NodeScorer interface {
	ScoreNode(
		node *tree.Node,
		siblings []*tree.Node,
	) float64
}

// MaximalSelector ...
type MaximalSelector struct {
	NodeScorer NodeScorer
}

// SelectNode ...
func (selector MaximalSelector) SelectNode(
	nodes []*tree.Node,
) *tree.Node {
	var maximum *tree.Node
	maximumScore := math.Inf(-1)
	for _, node := range nodes {
		nodeScore := selector.NodeScorer.
			ScoreNode(node, nodes)
		if nodeScore > maximumScore {
			maximum = node
			maximumScore = nodeScore
		}
	}

	return maximum
}
