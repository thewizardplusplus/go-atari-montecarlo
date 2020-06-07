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
	storage models.StoneStorage,
	color models.Color,
	settings searchSettings,
) (models.Move, error) {
	generator := models.MoveGenerator{}

	randomSelector := selectors.RandomMoveSelector{}
	generalSelector := selectors.MaximalNodeSelector{
		NodeScorer: scorers.UCBScorer{
			Factor: ucbFactor,
		},
	}

	var simulator simulators.Simulator // nolint: staticcheck
	simulator = simulators.RolloutSimulator{
		MoveGenerator: generator,
		MoveSelector:  randomSelector,
	}
	if settings.parallelSimulator {
		simulator = simulators.ParallelSimulator{
			Simulator:   simulator,
			Concurrency: runtime.NumCPU(),
		}
	}

	var bulkySimulator builders.BulkySimulator
	if !settings.parallelBulkySimulator {
		bulkySimulator = bulky.FirstNodeSimulator{
			Simulator: simulator,
		}
	} else {
		bulkySimulator = bulky.AllNodesSimulator{
			Simulator: simulator,
		}
	}

	var builder builders.Builder
	terminator := terminators.NewPassTerminator(settings.maximalPass)
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
		Move:    models.NewPreliminaryMove(color),
		Storage: storage,
	}
	searcher := searchers.MoveSearcher{
		MoveGenerator: generator,
		Builder:       builder,
		NodeSelector:  generalSelector,
	}
	node, err := searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}

func BenchmarkSearch_5Passes(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 5,
			},
		)
	}
}

func BenchmarkSearch_10Passes(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 10,
			},
		)
	}
}

func BenchmarkSearch_15Passes(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 15,
			},
		)
	}
}

func BenchmarkSearch_20Passes(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
		search(
			initialBoard,
			initialColor,
			searchSettings{
				maximalPass: 20,
			},
		)
	}
}

func BenchmarkSearch_5PassesAndParallelSimulator(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_10PassesAndParallelSimulator(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_15PassesAndParallelSimulator(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_20PassesAndParallelSimulator(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_5PassesAndParallelBulkySimulator(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_10PassesAndParallelBulkySimulator(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_15PassesAndParallelBulkySimulator(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_20PassesAndParallelBulkySimulator(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_5PassesAndParallelBuilder(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_10PassesAndParallelBuilder(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_15PassesAndParallelBuilder(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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

func BenchmarkSearch_20PassesAndParallelBuilder(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		// nolint: errcheck
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
