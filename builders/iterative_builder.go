package builders

import (
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// Builder ...
type Builder interface {
	Pass(root *tree.Node)
}

// IterativeBuilder ...
type IterativeBuilder struct {
	Builder    Builder
	PassCount  int
	Terminator terminators.BuildingTerminator
}

// Pass ...
func (builder IterativeBuilder) Pass(
	root *tree.Node,
) {
	for i := 0; i < builder.PassCount; i++ {
		builder.Builder.Pass(root)
	}
}
