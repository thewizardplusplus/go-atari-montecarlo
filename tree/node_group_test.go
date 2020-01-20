package tree

import (
	"testing"
)

func TestNodeGroupTotalGameCount(
	test *testing.T,
) {
	nodes := NodeGroup{
		&Node{
			State: NodeState{
				GameCount: 10,
				WinCount:  1,
			},
		},
		&Node{
			State: NodeState{
				GameCount: 10,
				WinCount:  2,
			},
		},
	}
	count := nodes.TotalGameCount()

	if count != 20 {
		test.Fail()
	}
}
