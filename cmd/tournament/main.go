package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

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
	defaultUCBFactor = 1
	defaultDuration  = 10 * time.Second
	gameCount        = 10
	startColor       = models.Black
)

var (
	startBoard = models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	lowVsHighUCBSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              0.1,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              10,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
	}
	usualVsParallelSimulatorSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			parallelSimulator:      true,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
	}
	usualVsParallelBulkySimulatorSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: true,
			parallelBuilder:        false,
		},
	}
	usualVsParallelBuilderSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        true,
		},
	}
)

type taskInbox chan func()

type searchingSettings struct {
	ucbFactor              float64
	maximalDuration        time.Duration
	parallelSimulator      bool
	parallelBulkySimulator bool
	parallelBuilder        bool
}

type integratedSearcher struct {
	terminator terminators.BuildingTerminator
	searcher   searchers.MoveSearcher
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
		terminators.NewTimeTerminator(
			time.Now,
			settings.maximalDuration,
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
	return integratedSearcher{
		terminator: terminator,
		searcher:   baseSearcher,
	}
}

func (searcher integratedSearcher) search(
	board models.Board,
	previousMove models.Move,
) (models.Move, error) {
	searcher.terminator.Reset()

	root := &tree.Node{
		Move:  previousMove,
		Board: board,
	}
	node, err :=
		searcher.searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}

type gameSettings struct {
	firstSearcher  searchingSettings
	secondSearcher searchingSettings
}

type score struct {
	winCount uint64
}

func (score *score) wins() uint64 {
	return atomic.LoadUint64(&score.winCount)
}

func (score *score) elo(
	gameCount uint64,
) float64 {
	winPercent := float64(score.wins()) /
		float64(gameCount)
	return 400 *
		math.Log10(winPercent/(1-winPercent))
}

func (score *score) win() {
	atomic.AddUint64(&score.winCount, 1)
}

type scores struct {
	firstSearcher  score
	secondSearcher score
}

func (scores *scores) games() uint64 {
	return scores.firstSearcher.wins() +
		scores.secondSearcher.wins()
}

func (scores *scores) addGame(
	errColor models.Color,
	err error,
) {
	switch err {
	case models.ErrAlreadyWin:
	case models.ErrAlreadyLoss:
		errColor = errColor.Negative()
	default:
		return
	}

	if errColor == startColor {
		scores.firstSearcher.win()
	} else {
		scores.secondSearcher.win()
	}
}

func (scores *scores) String() string {
	games := scores.games()
	return fmt.Sprintf(
		"Games: %d\n"+
			"First Searcher: %d\n"+
			"Second Searcher: %d\n"+
			"Second Searcher Elo Delta: %.2f",
		games,
		scores.firstSearcher.wins(),
		scores.secondSearcher.wins(),
		scores.secondSearcher.elo(games),
	)
}

func pool() (tasks taskInbox, wait func()) {
	threadCount := runtime.NumCPU()

	var waiter sync.WaitGroup
	waiter.Add(threadCount)

	tasks = make(taskInbox)
	for i := 0; i < threadCount; i++ {
		go func() {
			defer waiter.Done()
			fmt.Print("#")

			for task := range tasks {
				fmt.Print("%")
				task()
			}
		}()
	}

	return tasks, func() { waiter.Wait() }
}

func game(
	board models.Board,
	color models.Color,
	settings gameSettings,
) (models.Color, error) {
	firstSearcher := newIntegratedSearcher(
		settings.firstSearcher,
	)
	secondSearcher := newIntegratedSearcher(
		settings.secondSearcher,
	)

	previousMove :=
		models.NewPreliminaryMove(color)
	for ply := 0; ; ply++ {
		if ply%5 == 0 {
			fmt.Print(".")
		}

		var move models.Move
		var err error
		if ply%2 == 0 {
			move, err = firstSearcher.search(
				board,
				previousMove,
			)
		} else {
			move, err = secondSearcher.search(
				board,
				previousMove,
			)
		}
		if err != nil {
			errColor :=
				previousMove.Color.Negative()
			return errColor, err
		}

		board = board.ApplyMove(move)
		previousMove = move
	}
}

func markWinner(
	errColor models.Color,
	err error,
) {
	switch err {
	case models.ErrAlreadyWin:
	case models.ErrAlreadyLoss:
		errColor = errColor.Negative()
	default:
		log.Println(err)
		return
	}

	if errColor == startColor {
		fmt.Print("F")
	} else {
		fmt.Print("S")
	}
}

func main() {
	//settings := lowVsHighUCBSettings
	//settings := usualVsParallelSimulatorSettings
	//settings := usualVsParallelBulkySimulatorSettings
	settings := usualVsParallelBuilderSettings

	var scores scores
	tasks, wait := pool()
	for i := 0; i < gameCount; i++ {
		tasks <- func() {
			errColor, err := game(
				startBoard,
				startColor,
				settings,
			)
			scores.addGame(errColor, err)
			markWinner(errColor, err)
		}
	}

	close(tasks)
	wait()

	fmt.Println()
	fmt.Println(scores.String())
}
