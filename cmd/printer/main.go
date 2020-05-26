package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
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
	ucbFactor       = math.Sqrt2
	maximalDuration = 10 * time.Second
	minimalAxisCode = 97
	initialColor    = models.Black
)

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
		Builder:      builder,
		NodeSelector: generalSelector,
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
	return string(axis + minimalAxisCode)
}

func game(
	board models.Board,
	color models.Color,
	settings gameSettings,
) (history, models.Color, error) {
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
	cpuProf, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(
			"unable to create a CPU profile: ",
			err,
		)
	}
	defer cpuProf.Close()

	memProf, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(
			"unable to create a memory profile: ",
			err,
		)
	}
	defer memProf.Close()

	err = pprof.StartCPUProfile(cpuProf)
	if err != nil {
		log.Fatal(
			"unable to start CPU profiling: ",
			err,
		)
	}
	defer pprof.StopCPUProfile()

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

	history, errColor, err := game(
		initialBoard,
		initialColor,
		settings,
	)
	markWinner(errColor, err)
	fmt.Println(history)

	runtime.GC() // get up-to-date statistics

	err = pprof.WriteHeapProfile(memProf)
	if err != nil {
		log.Fatal(
			"unable to write a memory profile: ",
			err,
		)
	}
}
