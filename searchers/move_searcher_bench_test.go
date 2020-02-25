package searchers_test

import (
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
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
		search(board, models.Black, 1, 10)
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
		search(board, models.Black, 1, 100)
	}
}

func search(
	board models.Board,
	color models.Color,
	ucbFactor float64,
	maximalPass int,
) (models.Move, error) {
	root := tree.NewNode(board, color)
	randomSelector :=
		selectors.RandomSelector{}
	maximalSelector :=
		selectors.MaximalSelector{
			NodeScorer: scorers.UCBScorer{
				Factor: ucbFactor,
			},
		}
	simulator := simulators.RolloutSimulator{
		MoveSelector: selectors.MoveSelector{
			NodeSelector: randomSelector,
		},
	}
	terminator :=
		terminators.NewPassTerminator(
			maximalPass,
		)
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector: maximalSelector,
			Simulator:    simulator,
		},
		Terminator: terminator,
	}
	searcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: maximalSelector,
	}
	node, err := searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}
