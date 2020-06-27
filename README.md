# go-atari-montecarlo

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-atari-montecarlo?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-atari-montecarlo)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-atari-montecarlo)](https://goreportcard.com/report/github.com/thewizardplusplus/go-atari-montecarlo)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-atari-montecarlo.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-atari-montecarlo)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-atari-montecarlo/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-atari-montecarlo)

The library that implements an [Atari Go](https://senseis.xmp.net/?AtariGo) engine based on the Monte Carlo tree search algorithm.

_**Disclaimer:** this library was written directly on an Android smartphone with the AnGoIde IDE._

## Features

- move searching via the [Monte Carlo tree search algorithm](https://en.wikipedia.org/wiki/Monte_Carlo_tree_search):
  - move selectors:
    - random selecting;
    - selecting by a maximal node score:
      - scoring by the [Upper Confidence Bound algorithm](https://en.wikipedia.org/wiki/Multi-armed_bandit);
  - game simulating by simple random rollout;
  - tree building:
    - by a single pass;
    - by iterative passes:
      - iteration terminating:
        - by a pass;
        - by a time;
  - move searchers:
    - searcher that doesn't reuse a built tree;
- optimization via parallel move searching:
  - parallel game simulating:
    - of a single node child;
    - of all node children;
  - parallel tree building;
- easily extensible and composable architecture:
  - of move selectors:
    - of node scorers;
  - of game simulators;
  - of tree builders:
    - of iteration terminators;
  - of move searchers.

## Installation

```
$ go get github.com/thewizardplusplus/go-atari-montecarlo
```

## Examples

`searchers.MoveSearcher.SearchMove()` without parallelism:

```go
package main

import (
	"fmt"
	"log"
	"math"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators/bulky"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func main() {
	// +-+-+-+-+-+
	// |W|W|W|W|X|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	points := board.Size().Points()
	for _, point := range points[:len(points)-1] {
		board = board.ApplyMove(models.Move{Color: models.White, Point: point})
	}

	generator := models.MoveGenerator{}
	randomSelector := selectors.RandomMoveSelector{}
	generalSelector := selectors.MaximalNodeSelector{
		NodeScorer: scorers.UCBScorer{Factor: math.Sqrt2},
	}
	simulator := bulky.FirstNodeSimulator{
		Simulator: simulators.RolloutSimulator{
			MoveGenerator: generator,
			MoveSelector:  randomSelector,
		},
	}
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector:  generalSelector,
			MoveGenerator: generator,
			Simulator:     simulator,
		},
		Terminator: terminators.NewPassTerminator(2),
	}
	searcher := searchers.MoveSearcher{
		MoveGenerator: generator,
		Builder:       builder,
		NodeSelector:  generalSelector,
	}

	preliminaryMove := models.NewPreliminaryMove(models.Black)
	root := &tree.Node{Move: preliminaryMove, Storage: board}
	node, err := searcher.SearchMove(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", node.Move)

	// Output: {Color:0 Point:{Column:4 Row:4}}
}
```

`searchers.MoveSearcher.SearchMove()` with parallel game simulating of a single node child:

```go
package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators/bulky"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func main() {
	// +-+-+-+-+-+
	// |W|W|W|W|X|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	points := board.Size().Points()
	for _, point := range points[:len(points)-1] {
		board = board.ApplyMove(models.Move{Color: models.White, Point: point})
	}

	generator := models.MoveGenerator{}
	randomSelector := selectors.RandomMoveSelector{}
	generalSelector := selectors.MaximalNodeSelector{
		NodeScorer: scorers.UCBScorer{Factor: math.Sqrt2},
	}
	simulator := bulky.FirstNodeSimulator{
		Simulator: simulators.ParallelSimulator{
			Simulator: simulators.RolloutSimulator{
				MoveGenerator: generator,
				MoveSelector:  randomSelector,
			},
			Concurrency: runtime.NumCPU(),
		},
	}
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector:  generalSelector,
			MoveGenerator: generator,
			Simulator:     simulator,
		},
		Terminator: terminators.NewPassTerminator(2),
	}
	searcher := searchers.MoveSearcher{
		MoveGenerator: generator,
		Builder:       builder,
		NodeSelector:  generalSelector,
	}

	preliminaryMove := models.NewPreliminaryMove(models.Black)
	root := &tree.Node{Move: preliminaryMove, Storage: board}
	node, err := searcher.SearchMove(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", node.Move)

	// Output: {Color:0 Point:{Column:4 Row:4}}
}
```

`searchers.MoveSearcher.SearchMove()` with parallel game simulating of all node children:

```go
package main

import (
	"fmt"
	"log"
	"math"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators/bulky"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func main() {
	// +-+-+-+-+-+
	// |W|W|W|W|X|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	points := board.Size().Points()
	for _, point := range points[:len(points)-1] {
		board = board.ApplyMove(models.Move{Color: models.White, Point: point})
	}

	generator := models.MoveGenerator{}
	randomSelector := selectors.RandomMoveSelector{}
	generalSelector := selectors.MaximalNodeSelector{
		NodeScorer: scorers.UCBScorer{Factor: math.Sqrt2},
	}
	simulator := bulky.AllNodesSimulator{
		Simulator: simulators.RolloutSimulator{
			MoveGenerator: generator,
			MoveSelector:  randomSelector,
		},
	}
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector:  generalSelector,
			MoveGenerator: generator,
			Simulator:     simulator,
		},
		Terminator: terminators.NewPassTerminator(2),
	}
	searcher := searchers.MoveSearcher{
		MoveGenerator: generator,
		Builder:       builder,
		NodeSelector:  generalSelector,
	}

	preliminaryMove := models.NewPreliminaryMove(models.Black)
	root := &tree.Node{Move: preliminaryMove, Storage: board}
	node, err := searcher.SearchMove(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", node.Move)

	// Output: {Color:0 Point:{Column:4 Row:4}}
}
```

`searchers.MoveSearcher.SearchMove()` with parallel tree building:

```go
package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators/bulky"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func main() {
	// +-+-+-+-+-+
	// |W|W|W|W|X|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	// |W|W|W|W|W|
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	points := board.Size().Points()
	for _, point := range points[:len(points)-1] {
		board = board.ApplyMove(models.Move{Color: models.White, Point: point})
	}

	generator := models.MoveGenerator{}
	randomSelector := selectors.RandomMoveSelector{}
	generalSelector := selectors.MaximalNodeSelector{
		NodeScorer: scorers.UCBScorer{Factor: math.Sqrt2},
	}
	simulator := bulky.FirstNodeSimulator{
		Simulator: simulators.RolloutSimulator{
			MoveGenerator: generator,
			MoveSelector:  randomSelector,
		},
	}
	builder := builders.ParallelBuilder{
		Builder: builders.IterativeBuilder{
			Builder: builders.TreeBuilder{
				NodeSelector:  generalSelector,
				MoveGenerator: generator,
				Simulator:     simulator,
			},
			Terminator: terminators.NewPassTerminator(2),
		},
		Concurrency: runtime.NumCPU(),
	}
	searcher := searchers.MoveSearcher{
		MoveGenerator: generator,
		Builder:       builder,
		NodeSelector:  generalSelector,
	}

	preliminaryMove := models.NewPreliminaryMove(models.Black)
	root := &tree.Node{Move: preliminaryMove, Storage: board}
	node, err := searcher.SearchMove(root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", node.Move)

	// Output: {Color:0 Point:{Column:4 Row:4}}
}
```

## Benchmarks

Without parallelism:

```
BenchmarkSearch_5Passes-8                             	     300	   6041138 ns/op
BenchmarkSearch_10Passes-8                            	     100	  11545956 ns/op
BenchmarkSearch_15Passes-8                            	     100	  16472837 ns/op
BenchmarkSearch_20Passes-8                            	     100	  21192349 ns/op
```

With parallel game simulating of a single node child:

```
BenchmarkSearch_5PassesAndParallelSimulator-8         	     100	  11482508 ns/op
BenchmarkSearch_10PassesAndParallelSimulator-8        	      50	  24193690 ns/op
BenchmarkSearch_15PassesAndParallelSimulator-8        	      30	  35422514 ns/op
BenchmarkSearch_20PassesAndParallelSimulator-8        	      30	  46065076 ns/op
```

With parallel game simulating of all node children:

```
BenchmarkSearch_5PassesAndParallelBulkySimulator-8    	      50	  25558541 ns/op
BenchmarkSearch_10PassesAndParallelBulkySimulator-8   	      20	  59455409 ns/op
BenchmarkSearch_15PassesAndParallelBulkySimulator-8   	      20	  85564306 ns/op
BenchmarkSearch_20PassesAndParallelBulkySimulator-8   	      10	 186877684 ns/op
```

With parallel tree building:

```
BenchmarkSearch_5PassesAndParallelBuilder-8           	     100	  15637022 ns/op
BenchmarkSearch_10PassesAndParallelBuilder-8          	      50	  32202269 ns/op
BenchmarkSearch_15PassesAndParallelBuilder-8          	      30	  52235995 ns/op
BenchmarkSearch_20PassesAndParallelBuilder-8          	      20	  71661751 ns/op
```

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
