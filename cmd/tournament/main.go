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
	initialColor    = models.Black
	ucbFactor       = math.Sqrt2
	maximalDuration = 10 * time.Second
	gameCount       = 10
)

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

	if errColor == initialColor {
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

type taskInbox chan func()

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

type searchSettings struct {
	parallelSimulator      bool
	parallelBulkySimulator bool
	parallelBuilder        bool
}

func search(
	board models.Board,
	previousMove models.Move,
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
		terminators.NewTimeTerminator(
			time.Now,
			maximalDuration,
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
		Move:  previousMove,
		Board: board,
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

type gameSettings struct {
	firstSearcher  searchSettings
	secondSearcher searchSettings
}

func game(
	board models.Board,
	color models.Color,
	settings gameSettings,
) (models.Color, error) {
	previousMove :=
		models.NewPreliminaryMove(color)
	for ply := 0; ; ply++ {
		if ply%5 == 0 {
			fmt.Print(".")
		}

		var move models.Move
		var err error
		if ply%2 == 0 {
			move, err = search(
				board,
				previousMove,
				settings.firstSearcher,
			)
		} else {
			move, err = search(
				board,
				previousMove,
				settings.secondSearcher,
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

	if errColor == initialColor {
		fmt.Print("F")
	} else {
		fmt.Print("S")
	}
}

func main() {
	initialBoard := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)
	settings := gameSettings{
		firstSearcher: searchSettings{
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchSettings{
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        true,
		},
	}

	var scores scores
	tasks, wait := pool()
	for i := 0; i < gameCount; i++ {
		tasks <- func() {
			errColor, err := game(
				initialBoard,
				initialColor,
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
