package selectors

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockNodeSelector struct {
	selectNode func(
		nodes []*tree.Node,
	) *tree.Node
}

func (selector MockNodeSelector) SelectNode(
	nodes []*tree.Node,
) *tree.Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}

func TestMoveSelectorSelectMove(
	test *testing.T,
) {
	type fields struct {
		nodeSelector NodeSelector
	}
	type args struct {
		moves []models.Move
	}
	type data struct {
		fields fields
		args   args
		want   models.Move
	}

	for _, data := range []data{
		data{
			fields: fields{
				nodeSelector: MockNodeSelector{
					selectNode: func(
						nodes []*tree.Node,
					) *tree.Node {
						expectedNodes := []*tree.Node{
							&tree.Node{
								Move: models.Move{
									Color: models.White,
									Point: models.Point{
										Column: 1,
										Row:    2,
									},
								},
							},
							&tree.Node{
								Move: models.Move{
									Color: models.White,
									Point: models.Point{
										Column: 3,
										Row:    4,
									},
								},
							},
							&tree.Node{
								Move: models.Move{
									Color: models.White,
									Point: models.Point{
										Column: 5,
										Row:    6,
									},
								},
							},
						}
						if !reflect.DeepEqual(
							nodes,
							expectedNodes,
						) {
							test.Fail()
						}

						return nodes[1]
					},
				},
			},
			args: args{
				moves: []models.Move{
					models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					},
					models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 3,
							Row:    4,
						},
					},
					models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 5,
							Row:    6,
						},
					},
				},
			},
			want: models.Move{
				Color: models.White,
				Point: models.Point{
					Column: 3,
					Row:    4,
				},
			},
		},
	} {
		selector := MoveSelector{
			NodeSelector: data.fields.
				nodeSelector,
		}
		got := selector.
			SelectMove(data.args.moves)

		if !reflect.DeepEqual(got, data.want) {
			test.Fail()
		}
	}
}
