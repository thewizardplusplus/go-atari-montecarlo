package simulators

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockSimulator struct {
	simulate func(
		root *tree.Node,
	) tree.NodeState
}

func (simulator MockSimulator) Simulate(
	root *tree.Node,
) tree.NodeState {
	if simulator.simulate == nil {
		panic("not implemented")
	}

	return simulator.simulate(root)
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
			root *tree.Node,
		) tree.NodeState {
			if root.Move.Color != models.White {
				test.Fail()
			}

			expectedBoard := models.NewBoard(
				models.Size{
					Width:  3,
					Height: 3,
				},
			)
			if !reflect.DeepEqual(
				root.Board,
				expectedBoard,
			) {
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
	gotState := simulator.Simulate(&tree.Node{
		Move: models.Move{
			Color: models.White,
		},
		Board: board,
	})

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
