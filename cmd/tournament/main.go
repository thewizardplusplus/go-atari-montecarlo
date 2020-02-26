package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

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
		firstSearchingSettings: searchingSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
		secondSearchingSettings: searchingSettings{
			selectorType: ucbSelector,
			ucbFactor:    1,
			maximalPass:  10,
		},
		reuseSearchingTree: true,
	}
)

type taskInbox chan func()

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

type gameSettings struct {
	firstSearchingSettings  searchingSettings
	secondSearchingSettings searchingSettings
	reuseSearchingTree      bool
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
	gameSettings gameSettings,
) (models.Color, error) {
	for ply := 0; ; ply++ {
		if ply%5 == 0 {
			fmt.Print(".")
		}

		var searchingSettings searchingSettings
		if ply%2 == 0 {
			searchingSettings =
				gameSettings.firstSearchingSettings
		} else {
			searchingSettings =
				gameSettings.secondSearchingSettings
		}

		node, err :=
			search(root, searchingSettings)
		if err != nil {
			errColor := root.Move.Color.Negative()
			return errColor, err
		}

		if gameSettings.reuseSearchingTree {
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
	settings searchingSettings,
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

	tasks, wait := pool()
	for i := 0; i < gameCount; i++ {
		tasks <- func() {
			root :=
				tree.NewNode(board, firstColor)
			errColor, err := game(root, settings)
			markWinner(errColor, err)
		}
	}

	close(tasks)
	wait()

	fmt.Println()
	fmt.Println("done")
}
