package searchers_test

import (
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func BenchmarkSearch_10PassCount(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(board, models.Black, 10)
	}
}

func BenchmarkSearch_100PassCount(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(board, models.Black, 100)
	}
}

func BenchmarkSearch_1000PassCount(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(board, models.Black, 1000)
	}
}

func BenchmarkSearch_10000PassCount(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(board, models.Black, 10000)
	}
}

func search(
	board models.Board,
	color models.Color,
	passCount int,
) (move models.Move, ok bool) {
	root := tree.NewNode(board, color)
	var randomSelector selectors.RandomSelector
	maximalSelector :=
		selectors.MaximalSelector{
			NodeScorer: scorers.UCBScorer{
				Factor: 1,
			},
		}
	simulator := simulators.RolloutSimulator{
		MoveSelector: selectors.MoveSelector{
			NodeSelector: randomSelector,
		},
	}
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector: maximalSelector,
			Simulator:    simulator,
		},
		PassCount: passCount,
	}
	searcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: maximalSelector,
	}
	node, ok := searcher.SearchMove(root)
	if !ok {
		return models.Move{}, false
	}

	return node.Move, true
}
