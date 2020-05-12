package builders

import (
	syncutils "github.com/thewizardplusplus/go-atari-montecarlo/sync-utils"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// ParallelBuilder ...
type ParallelBuilder struct {
	Builder     Builder
	Concurrency int
}

// Pass ...
func (builder ParallelBuilder) Pass(
	root *tree.Node,
) {
	roots := syncutils.ParallelRun(
		builder.Concurrency,
		func(index int) (result interface{}) {
			rootCopy := root.ShallowCopy()
			builder.Builder.Pass(rootCopy)

			return rootCopy
		},
	)

	for _, rootCopy := range roots {
		root.MergeChildren(rootCopy.(*tree.Node))
	}
}
