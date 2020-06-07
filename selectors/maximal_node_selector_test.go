package selectors

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type WinRateNodeScorer struct{}

func (scorer WinRateNodeScorer) ScoreNode(node *tree.Node) float64 {
	return node.State.WinRate()
}

func TestMaximalNodeSelectorSelectNode(test *testing.T) {
	type fields struct {
		nodeScorer NodeScorer
	}
	type args struct {
		nodes tree.NodeGroup
	}
	type data struct {
		fields fields
		args   args
		want   *tree.Node
	}

	for _, data := range []data{
		{
			fields: fields{
				nodeScorer: WinRateNodeScorer{},
			},
			args: args{
				nodes: tree.NodeGroup{
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  1,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 3,
								Row:    4,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  3,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 5,
								Row:    6,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  5,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 7,
								Row:    8,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  4,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 9,
								Row:    10,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  2,
						},
					},
				},
			},
			want: &tree.Node{
				Move: models.Move{
					Point: models.Point{
						Column: 5,
						Row:    6,
					},
				},
				State: tree.NodeState{
					GameCount: 10,
					WinCount:  5,
				},
			},
		},
		{
			fields: fields{
				nodeScorer: WinRateNodeScorer{},
			},
			args: args{
				nodes: tree.NodeGroup{
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  1,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 3,
								Row:    4,
							},
						},
						State: tree.NodeState{
							GameCount: 10,
							WinCount:  3,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 5,
								Row:    6,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 7,
								Row:    8,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 9,
								Row:    10,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
				},
			},
			want: &tree.Node{
				Move: models.Move{
					Point: models.Point{
						Column: 5,
						Row:    6,
					},
				},
				State: tree.NodeState{
					GameCount: 0,
					WinCount:  0,
				},
			},
		},
		{
			fields: fields{
				nodeScorer: WinRateNodeScorer{},
			},
			args: args{
				nodes: tree.NodeGroup{
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 3,
								Row:    4,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 5,
								Row:    6,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 7,
								Row:    8,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
					&tree.Node{
						Move: models.Move{
							Point: models.Point{
								Column: 9,
								Row:    10,
							},
						},
						State: tree.NodeState{
							GameCount: 0,
							WinCount:  0,
						},
					},
				},
			},
			want: &tree.Node{
				Move: models.Move{
					Point: models.Point{
						Column: 1,
						Row:    2,
					},
				},
				State: tree.NodeState{
					GameCount: 0,
					WinCount:  0,
				},
			},
		},
	} {
		selector := MaximalNodeSelector{
			NodeScorer: data.fields.nodeScorer,
		}
		got := selector.SelectNode(data.args.nodes)

		if !reflect.DeepEqual(got, data.want) {
			test.Fail()
		}
	}
}
