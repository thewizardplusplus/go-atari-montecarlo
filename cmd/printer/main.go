package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
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
			reuseTree:              false,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              10,
			maximalDuration:        defaultDuration,
			reuseTree:              false,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
	}
	uniqueVsReusedTreeSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              false,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
	}
	reusedTreeVsParallelSimulatorSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      true,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
	}
	reusedTreeVsParallelBulkySimulatorSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: true,
			parallelBuilder:        false,
		},
	}
	reusedTreeVsParallelBuilderSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        true,
		},
	}
	reusedTreeVsAllParallelSettings = gameSettings{
		firstSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			ucbFactor:              defaultUCBFactor,
			maximalDuration:        defaultDuration,
			reuseTree:              true,
			parallelSimulator:      true,
			parallelBulkySimulator: true,
			parallelBuilder:        true,
		},
	}
)

type searchingSettings struct {
	ucbFactor              float64
	maximalDuration        time.Duration
	reuseTree              bool
	parallelSimulator      bool
	parallelBulkySimulator bool
	parallelBuilder        bool
}

type integratedSearcher struct {
	terminator terminators.BuildingTerminator
	searcher   searchers.Searcher
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

	return searcher
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

func (history history) String() string {
	commands := []string{";FF[4]"}
	for _, move := range history {
		command := newMoveCommand(move)
		commands = append(commands, command)
	}

	result := strings.Join(commands, "")
	return fmt.Sprintf("(%s)", result)
}

func newMoveCommand(
	move models.Move,
) string {
	var color string
	switch move.Color {
	case models.Black:
		color = "B"
	case models.White:
		color = "W"
	}

	column :=
		convertPointAxis(move.Point.Column)
	row := convertPointAxis(move.Point.Row)

	return fmt.Sprintf(
		";%s[%s%s]",
		color,
		column,
		row,
	)
}

func convertPointAxis(axis int) string {
	return string(axis + 97)
}

func game(
	board models.Board,
	color models.Color,
	settings gameSettings,
) (history, models.Color, error) {
	firstSearcher := newIntegratedSearcher(
		settings.firstSearcher,
	)
	secondSearcher := newIntegratedSearcher(
		settings.secondSearcher,
	)

	var history history
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
	//settings := lowVsHighUCBSettings
	//settings := uniqueVsReusedTreeSettings
	//settings := reusedTreeVsParallelSimulatorSettings
	//settings := reusedTreeVsParallelBulkySimulatorSettings
	settings := reusedTreeVsParallelBuilderSettings
	//settings := reusedTreeVsAllParallelSettings

	history, errColor, err := game(
		startBoard,
		startColor,
		settings,
	)
	markWinner(errColor, err)
	fmt.Println(history)
}
