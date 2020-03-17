package simulators

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockSimulator struct {
	simulate func(
		board models.Board,
		color models.Color,
	) tree.NodeState
}

func (simulator MockSimulator) Simulate(
	board models.Board,
	color models.Color,
) tree.NodeState {
	if simulator.simulate == nil {
		panic("not implemented")
	}

	return simulator.simulate(board, color)
}

func TestParallelSimulatorSimulate(
	test *testing.T,
) {
	board := models.NewBoard(
		models.Size{
			Width:  3,
			Height: 3,
		},
	)
	innerSimulator := MockSimulator{
		simulate: func(
			board models.Board,
			color models.Color,
		) tree.NodeState {
			expectedBoard := models.NewBoard(
				models.Size{
					Width:  3,
					Height: 3,
				},
			)
			if !reflect.DeepEqual(
				board,
				expectedBoard,
			) {
				test.Fail()
			}

			if color != models.White {
				test.Fail()
			}

			return tree.NodeState{
				GameCount: 3,
				WinCount:  2,
			}
		},
	}
	simulator := ParallelSimulator{
		Simulator:   innerSimulator,
		Concurrency: 10,
	}
	gotState :=
		simulator.Simulate(board, models.White)

	wantState := tree.NodeState{
		GameCount: 30,
		WinCount:  20,
	}
	if !reflect.DeepEqual(
		gotState,
		wantState,
	) {
		test.Fail()
	}
}
