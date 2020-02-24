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
	Terminator terminators.BuildingTerminator
}

// Pass ...
func (builder IterativeBuilder) Pass(
	root *tree.Node,
) {
	for pass := 0; !builder.Terminator.
		IsBuildingTerminated(pass); pass++ {
		builder.Builder.Pass(root)
	}
}
