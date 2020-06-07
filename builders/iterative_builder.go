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
func (builder IterativeBuilder) Pass(root *tree.Node) {
	isBuildingTerminated := builder.Terminator.IsBuildingTerminated
	for pass := 0; !isBuildingTerminated(pass); pass++ {
		builder.Builder.Pass(root)
	}
}
