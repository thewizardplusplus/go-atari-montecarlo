package builders

import (
	"sync"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// BuilderFactory ...
type BuilderFactory func() Builder

// ParallelBuilder ...
type ParallelBuilder struct {
	BuilderFactory BuilderFactory
	Concurrency    int
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
			innerBuilder :=
				builder.BuilderFactory()
			innerBuilder.Pass(rootCopy)

			roots <- rootCopy
		}()
	}

	waiter.Wait()
	close(roots)

	root.Children = (<-roots).Children
	for rootCopy := range roots {
		root.Children.Merge(rootCopy.Children)
	}
}
