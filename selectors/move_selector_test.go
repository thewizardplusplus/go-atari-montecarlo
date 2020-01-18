package selectors

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockNodeSelector struct {
	selectNode func(
		nodes []*tree.Node,
	) *tree.Node
}

func (selector MockNodeSelector) SelectNode(
	nodes []*tree.Node,
) *tree.Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}
