package scorers

import (
	"math"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestWinRateScorerScoreNode(
	test *testing.T,
) {
	type args struct {
		node *tree.Node
	}
	type data struct {
		args args
		want float64
	}

	for _, data := range []data{
		data{
			args: args{
				node: &tree.Node{
					State: tree.NodeState{
						GameCount: 4,
						WinCount:  2,
					},
				},
			},
			want: 0.5,
		},
		data{
			args: args{
				node: &tree.Node{
					State: tree.NodeState{
						GameCount: 0,
						WinCount:  0,
					},
				},
			},
			want: math.Inf(+1),
		},
	} {
		var scorer WinRateScorer
		got := scorer.ScoreNode(data.args.node)

		if got != data.want {
			test.Fail()
		}
	}
}
