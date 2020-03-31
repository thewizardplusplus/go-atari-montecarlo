package searchers_test

import (
	"runtime"
	"testing"

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

type searchingSettings struct {
	ucbFactor              float64
	maximalPass            int
	reuseTree              bool
	parallelSimulator      bool
	parallelBulkySimulator bool
	parallelBuilder        bool
}

type integratedSearcher struct {
	searcher searchers.Searcher
}

func newIntegratedSearcher(
	settings searchingSettings,
) integratedSearcher {
	generalSelector :=
		selectors.MaximalSelector{
			NodeScorer: scorers.UCBScorer{
				Factor: settings.ucbFactor,
			},
		}

	var simulator simulators.Simulator
	simulator = simulators.RolloutSimulator{
		MoveSelector: selectors.MoveSelector{
			NodeSelector: selectors.
				RandomSelector{},
		},
	}
	if settings.parallelSimulator {
		simulator =
			simulators.ParallelSimulator{
				Simulator:   simulator,
				Concurrency: runtime.NumCPU(),
			}
	}

	var bulkySimulator builders.BulkySimulator
	if !settings.parallelBulkySimulator {
		bulkySimulator =
			bulky.FirstNodeSimulator{
				Simulator: simulator,
			}
	} else {
		bulkySimulator =
			bulky.AllNodesSimulator{
				Simulator: simulator,
			}
	}

	var builder builders.Builder
	terminator :=
		terminators.NewPassTerminator(
			settings.maximalPass,
		)
	builder = builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector: generalSelector,
			Simulator:    bulkySimulator,
		},
		Terminator: terminator,
	}
	if settings.parallelBuilder {
		builder = builders.ParallelBuilder{
			Builder:     builder,
			Concurrency: runtime.NumCPU(),
		}
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

	return searcher
}

func (searcher integratedSearcher) search(
	board models.Board,
	color models.Color,
) (models.Move, error) {
	root :=
		tree.NewPreliminaryNode(board, color)
	node, err :=
		searcher.searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}

func BenchmarkSearch_10Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            10,
			reuseTree:              false,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
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

func BenchmarkSearch_100Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            100,
			reuseTree:              false,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeAnd10Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            10,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeAnd100Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            100,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeParallelSimulatorAnd10Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            10,
			reuseTree:              true,
			parallelSimulator:      true,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeParallelSimulatorAnd100Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            100,
			reuseTree:              true,
			parallelSimulator:      true,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeParallelBulkySimulatorAnd10Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            10,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: true,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeParallelBulkySimulatorAnd100Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            100,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: true,
			parallelBuilder:        false,
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

func BenchmarkSearch_reusedTreeParallelBuilderAnd10Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            10,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        true,
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

func BenchmarkSearch_reusedTreeParallelBuilderAnd100Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            100,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        true,
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

func BenchmarkSearch_reusedTreeAllParallelAnd10Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            10,
			reuseTree:              true,
			parallelSimulator:      true,
			parallelBulkySimulator: true,
			parallelBuilder:        true,
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

func BenchmarkSearch_reusedTreeAllParallelAnd100Passes(
	benchmark *testing.B,
) {
	searcher := newIntegratedSearcher(
		searchingSettings{
			ucbFactor:              1,
			maximalPass:            100,
			reuseTree:              true,
			parallelSimulator:      true,
			parallelBulkySimulator: true,
			parallelBuilder:        true,
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
