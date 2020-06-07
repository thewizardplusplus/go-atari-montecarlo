package bulky

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestAllNodesSimulatorSimulate(test *testing.T) {
	innerSimulator := MockSimulator{
		simulate: func(root *tree.Node) tree.NodeState {
			return tree.NodeState{
				GameCount: root.State.GameCount * 2,
				WinCount:  root.State.WinCount * 2,
			}
		},
	}
	simulator := AllNodesSimulator{
		Simulator: innerSimulator,
	}
	gotStates := simulator.
		Simulate(
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
		{
			GameCount: 4,
			WinCount:  2,
		},
		{
			GameCount: 8,
			WinCount:  6,
		},
	}
	if !reflect.DeepEqual(gotStates, wantStates) {
		test.Fail()
	}
}
