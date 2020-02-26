package main

import (
	"errors"
	"fmt"

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

type gameSettings struct {
	firstSearchingSettings  searchingSettings
	secondSearchingSettings searchingSettings
	reuseSearchingTree      bool
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

func main() {
	searchingSettings := searchingSettings{
		selectorType: ucbSelector,
		ucbFactor:    1,
		maximalPass:  10,
	}
	gameSettings := gameSettings{
		firstSearchingSettings:  searchingSettings,
		secondSearchingSettings: searchingSettings,
		reuseSearchingTree:      true,
	}

	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)
	root := tree.NewNode(board, models.Black)
	errColor, err := game(root, gameSettings)

	fmt.Println()
	fmt.Println(errColor, err)
}
