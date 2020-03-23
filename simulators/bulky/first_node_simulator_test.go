package bulky

import (
	"reflect"
	"testing"

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

func TestFirstNodeSimulatorSimulate(
	test *testing.T,
) {
	innerSimulator := MockSimulator{
		simulate: func(
			root *tree.Node,
		) tree.NodeState {
			expectedRoot := &tree.Node{
				State: tree.NodeState{
					GameCount: 2,
					WinCount:  1,
				},
			}
			if !reflect.DeepEqual(
				root,
				expectedRoot,
			) {
				test.Fail()
			}

			return tree.NodeState{
				GameCount: 6,
				WinCount:  5,
			}
		},
	}
	simulator := FirstNodeSimulator{
		Simulator: innerSimulator,
	}
	gotStates := simulator.Simulate(
		tree.NodeGroup{
			&tree.Node{
				State: tree.NodeState{
					GameCount: 2,
					WinCount:  1,
				},
			},
			&tree.Node{
				State: tree.NodeState{
					GameCount: 4,
					WinCount:  3,
				},
			},
		},
	)

	wantStates := []tree.NodeState{
		tree.NodeState{
			GameCount: 6,
			WinCount:  5,
		},
	}
	if !reflect.DeepEqual(
		gotStates,
		wantStates,
	) {
		test.Fail()
	}
}
