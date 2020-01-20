package selectors

import (
	"math/rand"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// RandomSelector ...
type RandomSelector struct{}

// SelectNode ...
func (selector RandomSelector) SelectNode(
	nodes tree.NodeGroup,
) *tree.Node {
	index := rand.Intn(len(nodes))
	return nodes[index]
}
