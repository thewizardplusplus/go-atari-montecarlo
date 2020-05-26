package searchers_test

import (
	"math"
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

const (
	initialColor = models.Black
	ucbFactor    = math.Sqrt2
)

var (
	initialBoard = models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)
)

type searchSettings struct {
	maximalPass            int
	parallelSimulator      bool
	parallelBulkySimulator bool
	parallelBuilder        bool
}

func search(
	board models.Board,
	color models.Color,
	settings searchSettings,
) (models.Move, error) {
	generator := models.MoveGenerator{}
	randomSelector :=
		selectors.RandomMoveSelector{}
	generalSelector :=
		selectors.MaximalNodeSelector{
			NodeScorer: scorers.UCBScorer{
				Factor: ucbFactor,
			},
		}

	var simulator simulators.Simulator
	simulator = simulators.RolloutSimulator{
		MoveGenerator: generator,
		MoveSelector:  randomSelector,
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
			NodeSelector:  generalSelector,
			MoveGenerator: generator,
			Simulator:     bulkySimulator,
		},
		Terminator: terminator,
	}
	if settings.parallelBuilder {
		builder = builders.ParallelBuilder{
			Builder:     builder,
			Concurrency: runtime.NumCPU(),
		}
	}

	root := &tree.Node{
		Move:  models.NewPreliminaryMove(color),
		Board: board,
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

func BenchmarkSearch_with5Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 5,
			},
		)
	}
}

func BenchmarkSearch_with10Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 10,
			},
		)
	}
}

func BenchmarkSearch_with15Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 15,
			},
		)
	}
}

func BenchmarkSearch_with20Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 20,
			},
		)
	}
}

func BenchmarkSearch_withParallelSimulatorAnd5Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:       5,
				parallelSimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelSimulatorAnd10Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:       10,
				parallelSimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelSimulatorAnd15Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:       15,
				parallelSimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelSimulatorAnd20Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:       20,
				parallelSimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBulkySimulatorAnd5Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:            5,
				parallelBulkySimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBulkySimulatorAnd10Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:            10,
				parallelBulkySimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBulkySimulatorAnd15Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:            15,
				parallelBulkySimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBulkySimulatorAnd20Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:            20,
				parallelBulkySimulator: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBuilderAnd5Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:     5,
				parallelBuilder: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBuilderAnd10Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:     10,
				parallelBuilder: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBuilderAnd15Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:     15,
				parallelBuilder: true,
			},
		)
	}
}

func BenchmarkSearch_withParallelBuilderAnd20Passes(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass:     20,
				parallelBuilder: true,
			},
		)
	}
}
