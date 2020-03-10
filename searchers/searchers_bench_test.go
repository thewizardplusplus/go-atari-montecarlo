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
	reuseTree    bool
}

type integratedSearcher struct {
	searcher searchers.Searcher
}

func newIntegratedSearcher(
	settings searchingSettings,
) (integratedSearcher, error) {
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
		return integratedSearcher{},
			errors.New("unknown selector type")
	}

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
	baseSearcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: generalSelector,
	}

	var searcher integratedSearcher
	if !settings.reuseTree {
		searcher.searcher = baseSearcher
	} else {
		searcher.searcher =
			searchers.FallbackSearcher{
				PrimarySearcher: searchers.
					NewReusedSearcher(baseSearcher),
				FallbackSearcher: baseSearcher,
			}
	}

	return searcher, nil
}

func (searcher integratedSearcher) search(
	board models.Board,
	color models.Color,
) (models.Move, error) {
	root := tree.NewNode(board, color)
	node, err :=
		searcher.searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}

func BenchmarkSearch_randomSelectorAnd10Passes(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: randomSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_randomSelectorAnd100Passes(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: randomSelector,
			ucbFactor:    1,
			maximalPass:  100,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_winRateSelectorAnd10Passes(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: winRateSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_winRateSelectorAnd100Passes(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: winRateSelector,
			ucbFactor:    1,
			maximalPass:  100,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_ucbSelectorAnd10Passes(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_ucbSelectorAnd100Passes(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  100,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_randomSelectorAnd10PassesAndReusedTree(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: randomSelector,
			ucbFactor:    1,
			maximalPass:  10,
			reuseTree:    true,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_randomSelectorAnd100PassesAndReusedTree(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: randomSelector,
			ucbFactor:    1,
			maximalPass:  100,
			reuseTree:    true,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_winRateSelectorAnd10PassesAndReusedTree(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: winRateSelector,
			ucbFactor:    1,
			maximalPass:  10,
			reuseTree:    true,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_winRateSelectorAnd100PassesAndReusedTree(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: winRateSelector,
			ucbFactor:    1,
			maximalPass:  100,
			reuseTree:    true,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_ucbSelectorAnd10PassesAndReusedTree(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  10,
			reuseTree:    true,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}

func BenchmarkSearch_ucbSelectorAnd100PassesAndReusedTree(
	benchmark *testing.B,
) {
	searcher, _ := newIntegratedSearcher(
		searchingSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  100,
			reuseTree:    true,
		},
	)
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	for i := 0; i < benchmark.N; i++ {
		searcher.search(board, models.Black)
	}
}
