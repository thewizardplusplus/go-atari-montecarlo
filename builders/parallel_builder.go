package builders

import (
	"sync"

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
	roots := make(
		tree.NodeGroup,
		builder.Concurrency,
	)

	var waiter sync.WaitGroup
	for i := 0; i < len(roots); i++ {
		waiter.Add(1)

		go func(i int) {
			defer waiter.Done()

			roots[i] = root.ShallowCopy()
			builder.Builder.Pass(roots[i])
		}(i)
	}
	waiter.Wait()

	for _, rootCopy := range roots {
		root.MergeChildren(rootCopy)
	}
}
