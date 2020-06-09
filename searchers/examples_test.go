package searchers_test

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

func ExampleMoveSearcher_withoutParallelism() {
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

func ExampleMoveSearcher_withParallelSimulator() {
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

func ExampleMoveSearcher_withParallelBulkySimulator() {
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

func ExampleMoveSearcher_withParallelBuilder() {
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
