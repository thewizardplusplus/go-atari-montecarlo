package main

import (
	"errors"
	"fmt"
	"log"
	"time"

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
	defaultUCBFactor = 1
	defaultDuration  = 1 * time.Second
	startColor       = models.Black
)

var (
	startBoard = models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)

	randomVsWinRateSettings = gameSettings{
		firstSearcher: searchingSettings{
			selectorType:    randomSelector,
			ucbFactor:       defaultUCBFactor,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
		secondSearcher: searchingSettings{
			selectorType:    winRateSelector,
			ucbFactor:       defaultUCBFactor,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
	}
	winRateVsUCBSettings = gameSettings{
		firstSearcher: searchingSettings{
			selectorType:    winRateSelector,
			ucbFactor:       defaultUCBFactor,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
		secondSearcher: searchingSettings{
			selectorType:    ucbSelector,
			ucbFactor:       defaultUCBFactor,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
	}
	lowVsHighUCBSettings = gameSettings{
		firstSearcher: searchingSettings{
			selectorType:    ucbSelector,
			ucbFactor:       0.1,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
		secondSearcher: searchingSettings{
			selectorType:    ucbSelector,
			ucbFactor:       10,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
	}
	uniqueVsReusedTreeSettings = gameSettings{
		firstSearcher: searchingSettings{
			selectorType:    ucbSelector,
			ucbFactor:       defaultUCBFactor,
			maximalDuration: defaultDuration,
			reuseTree:       false,
		},
		secondSearcher: searchingSettings{
			selectorType:    ucbSelector,
			ucbFactor:       defaultUCBFactor,
			maximalDuration: defaultDuration,
			reuseTree:       true,
		},
	}
)

type selectorType int

const (
	randomSelector selectorType = iota
	winRateSelector
	ucbSelector
)

type searchingSettings struct {
	selectorType    selectorType
	ucbFactor       float64
	maximalDuration time.Duration
	reuseTree       bool
}

type integratedSearcher struct {
	terminator terminators.BuildingTerminator
	searcher   searchers.Searcher
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
		terminators.NewTimeTerminator(
			time.Now,
			settings.maximalDuration,
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

	searcher := integratedSearcher{
		terminator: terminator,
	}
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

type history []models.Move

func game(
	board models.Board,
	color models.Color,
	settings gameSettings,
) (history, models.Color, error) {
	firstSearcher, err :=
		newIntegratedSearcher(
			settings.firstSearcher,
		)
	if err != nil {
		return nil, 0, err
	}

	secondSearcher, err :=
		newIntegratedSearcher(
			settings.secondSearcher,
		)
	if err != nil {
		return nil, 0, err
	}

	var history history
	previousMove :=
		models.NewPreliminaryMove(color)
	for ply := 0; ; ply++ {
		if ply%5 == 0 {
			fmt.Print(".")
		}

		var move models.Move
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
			errColor := move.Color.Negative()
			return history, errColor, err
		}

		board = board.ApplyMove(move)
		history = append(history, move)
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
		fmt.Println("F")
	} else {
		fmt.Println("S")
	}
}

func main() {
	//settings := randomVsWinRateSettings
	//settings := winRateVsUCBSettings
	//settings := lowVsHighUCBSettings
	settings := uniqueVsReusedTreeSettings

	history, errColor, err := game(
		startBoard,
		startColor,
		settings,
	)
	markWinner(errColor, err)
	fmt.Println(history)
}
