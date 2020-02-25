package searchers_test

import (
	"errors"
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

type selectorType int

const (
	randomSelector selectorType = iota
	winRateSelector
	ucbSelector
)

type searchingSettings struct {
	selectorType selectorType
	ucbFactor    float64
	maximalPass  int
}

func BenchmarkSearch_randomSelectorAnd10Passes(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(
			board,
			models.Black,
			searchingSettings{
				selectorType: randomSelector,
				ucbFactor:    1,
				maximalPass:  10,
			},
		)
	}
}

func BenchmarkSearch_randomSelectorAnd100Passes(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(
			board,
			models.Black,
			searchingSettings{
				selectorType: randomSelector,
				ucbFactor:    1,
				maximalPass:  100,
			},
		)
	}
}

func BenchmarkSearch_winRateSelectorAnd10Passes(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(
			board,
			models.Black,
			searchingSettings{
				selectorType: winRateSelector,
				ucbFactor:    1,
				maximalPass:  10,
			},
		)
	}
}

func BenchmarkSearch_winRateSelectorAnd100Passes(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(
			board,
			models.Black,
			searchingSettings{
				selectorType: winRateSelector,
				ucbFactor:    1,
				maximalPass:  100,
			},
		)
	}
}

func BenchmarkSearch_ucbSelectorAnd10Passes(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(
			board,
			models.Black,
			searchingSettings{
				selectorType: ucbSelector,
				ucbFactor:    1,
				maximalPass:  10,
			},
		)
	}
}

func BenchmarkSearch_ucbSelectorAnd100Passes(
	benchmark *testing.B,
) {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		search(
			board,
			models.Black,
			searchingSettings{
				selectorType: ucbSelector,
				ucbFactor:    1,
				maximalPass:  100,
			},
		)
	}
}

func search(
	board models.Board,
	color models.Color,
	settings searchingSettings,
) (models.Move, error) {
	var generalSelector tree.NodeSelector
	switch settings.selectorType {
	case randomSelector:
		generalSelector =
			selectors.RandomSelector{}
	case winRateSelector:
		generalSelector =
			selectors.MaximalSelector{
				NodeScorer: scorers.WinRateScorer{},
			}
	case ucbSelector:
		generalSelector =
			selectors.MaximalSelector{
				NodeScorer: scorers.UCBScorer{
					Factor: settings.ucbFactor,
				},
			}
	default:
		return models.Move{},
			errors.New("unknown selector type")
	}

	root := tree.NewNode(board, color)
	randomSelector :=
		selectors.RandomSelector{}
	simulator := simulators.RolloutSimulator{
		MoveSelector: selectors.MoveSelector{
			NodeSelector: randomSelector,
		},
	}
	terminator :=
		terminators.NewPassTerminator(
			settings.maximalPass,
		)
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector: generalSelector,
			Simulator:    simulator,
		},
		Terminator: terminator,
	}
	searcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: generalSelector,
	}
	node, err := searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}
