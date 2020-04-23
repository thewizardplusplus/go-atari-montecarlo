package main

import (
	"fmt"
	"log"
	"math"
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
	initialColor    = models.Black
	ucbFactor       = math.Sqrt2
	maximalDuration = 10 * time.Second
)

type searchingSettings struct {
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
		MoveSelector: randomSelector,
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

	if errColor == initialColor {
		fmt.Println("F")
	} else {
		fmt.Println("S")
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
		firstSearcher: searchingSettings{
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        false,
		},
		secondSearcher: searchingSettings{
			parallelSimulator:      false,
			parallelBulkySimulator: false,
			parallelBuilder:        true,
		},
	}

	history, errColor, err := game(
		initialBoard,
		initialColor,
		settings,
	)
	markWinner(errColor, err)
	fmt.Println(history)
}
