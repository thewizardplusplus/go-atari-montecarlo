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
	var waiter sync.WaitGroup
	waiter.Add(builder.Concurrency)

	roots := make(
		chan *tree.Node,
		builder.Concurrency,
	)
	concurrency := builder.Concurrency
	for i := 0; i < concurrency; i++ {
		go func() {
			defer waiter.Done()

			rootCopy := root.ShallowCopy()
			builder.Builder.Pass(rootCopy)

			roots <- rootCopy
		}()
	}

	waiter.Wait()
	close(roots)

	for rootCopy := range roots {
		root.MergeChildren(rootCopy)
	}
}
