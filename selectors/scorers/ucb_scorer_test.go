package scorers

import (
	"math"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestUCBScorerScoreNode(
	test *testing.T,
) {
	type fields struct {
		factor float64
	}
	type args struct {
		node *tree.Node
	}
	type data struct {
		fields fields
		args   args
		want   float64
	}

	for _, data := range []data{
		data{
			fields: fields{
				factor: 2,
			},
			args: args{
				node: &tree.Node{
					Parent: &tree.Node{
						State: tree.NodeState{
							GameCount: 9,
							WinCount:  5,
						},
					},
					State: tree.NodeState{
						GameCount: 4,
						WinCount:  2,
					},
				},
			},
			want: 1.98,
		},
		data{
			fields: fields{
				factor: 2,
			},
			args: args{
				node: &tree.Node{
					Parent: &tree.Node{
						State: tree.NodeState{
							GameCount: 5,
							WinCount:  3,
						},
					},
					State: tree.NodeState{
						GameCount: 0,
						WinCount:  0,
					},
				},
			},
			want: math.Inf(+1),
		},
		data{
			fields: fields{
				factor: 2,
			},
			args: args{
				node: &tree.Node{
					Parent: &tree.Node{
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					State: tree.NodeState{
						GameCount: 0,
						WinCount:  0,
					},
				},
			},
			want: math.Inf(+1),
		},
	} {
		scorer := UCBScorer{
			Factor: data.fields.factor,
		}
		got := scorer.ScoreNode(data.args.node)

		roundedGot := math.Floor(got*100) / 100
		if roundedGot != data.want {
			test.Fail()
		}
	}
}
