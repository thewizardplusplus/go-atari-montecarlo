package main

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

const (
	gameCount  = 10
	firstColor = models.Black
)

var (
	settings = gameSettings{
		firstSearcher: searcherSettings{
			selectorType: randomSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
		secondSearcher: searcherSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
		reuseTree: false,
	}
)

type taskInbox chan func()

type selectorType int

const (
	randomSelector selectorType = iota
	winRateSelector
	ucbSelector
)

type searcherSettings struct {
	selectorType selectorType
	ucbFactor    float64
	maximalPass  int
}

type gameSettings struct {
	firstSearcher  searcherSettings
	secondSearcher searcherSettings
	reuseTree      bool
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

	if errColor == firstColor {
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
	root *tree.Node,
	settings gameSettings,
) (models.Color, error) {
	for ply := 0; ; ply++ {
		if ply%5 == 0 {
			fmt.Print(".")
		}

		var searcherSettings searcherSettings
		if ply%2 == 0 {
			searcherSettings =
				settings.firstSearcher
		} else {
			searcherSettings =
				settings.secondSearcher
		}

		node, err :=
			search(root, searcherSettings)
		if err != nil {
			errColor := root.Move.Color.Negative()
			return errColor, err
		}

		if settings.reuseTree {
			root = node
		} else {
			root = tree.NewNode(
				node.Board,
				node.Move.Color.Negative(),
			)
		}
	}
}

func search(
	root *tree.Node,
	settings searcherSettings,
) (*tree.Node, error) {
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
		return nil,
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
	searcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: generalSelector,
	}
	return searcher.SearchMove(root)
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
		fmt.Print("E")
	}

	if errColor == firstColor {
		fmt.Print("F")
	} else {
		fmt.Print("S")
	}
}

func main() {
	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	var scores scores
	tasks, wait := pool()
	for i := 0; i < gameCount; i++ {
		tasks <- func() {
			root :=
				tree.NewNode(board, firstColor)
			errColor, err := game(root, settings)
			scores.addGame(errColor, err)
			markWinner(errColor, err)
		}
	}

	close(tasks)
	wait()

	fmt.Println()
	fmt.Println(scores.String())
}
