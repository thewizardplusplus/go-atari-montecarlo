package searchers

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockBuilder struct {
	pass func(root *tree.Node)
}

func (builder MockBuilder) Pass(
	root *tree.Node,
) {
	if builder.pass == nil {
		panic("not implemented")
	}

	builder.pass(root)
}

type MockNodeSelector struct {
	selectNode func(
		nodes tree.NodeGroup,
	) *tree.Node
}

func (selector MockNodeSelector) SelectNode(
	nodes tree.NodeGroup,
) *tree.Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}

func TestMoveSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		moveGenerator models.Generator
		builder       builders.Builder
		nodeSelector  tree.NodeSelector
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields   fields
		args     args
		wantNode *tree.Node
		wantErr  error
	}

	for _, data := range []data{
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				builder: MockBuilder{
					pass: func(root *tree.Node) {
						panic("not implemented")
					},
				},
				nodeSelector: MockNodeSelector{
					selectNode: func(
						nodes tree.NodeGroup,
					) *tree.Node {
						panic("not implemented")
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    2,
						},
					},
					Board: func() models.Board {
						board := models.NewBoard(
							models.Size{
								Width:  3,
								Height: 3,
							},
						)

						points := board.Size().Points()
						for _, point := range points {
							move := models.Move{
								Color: models.White,
								Point: point,
							}
							board = board.ApplyMove(move)
						}

						return board
					}(),
					State: tree.NodeState{
						GameCount: 2,
						WinCount:  1,
					},
				},
			},
			wantNode: nil,
			wantErr:  models.ErrAlreadyWin,
		},
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				builder: MockBuilder{
					pass: func(root *tree.Node) {
						expectedRoot := &tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.NilPoint,
							},
							Board: func() models.Board {
								return models.NewBoard(
									models.Size{
										Width:  3,
										Height: 3,
									},
								)
							}(),
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

						// don't change the passed node
					},
				},
				nodeSelector: MockNodeSelector{
					selectNode: func(
						nodes tree.NodeGroup,
					) *tree.Node {
						panic("not implemented")
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.NilPoint,
					},
					Board: func() models.Board {
						return models.NewBoard(
							models.Size{
								Width:  3,
								Height: 3,
							},
						)
					}(),
					State: tree.NodeState{
						GameCount: 2,
						WinCount:  1,
					},
				},
			},
			wantNode: nil,
			wantErr:  ErrFailedBuilding,
		},
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				builder: MockBuilder{
					pass: func(root *tree.Node) {
						expectedRoot := &tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.NilPoint,
							},
							Board: func() models.Board {
								return models.NewBoard(
									models.Size{
										Width:  3,
										Height: 3,
									},
								)
							}(),
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

						// add few children
						// to the passed node
						root.Children = tree.NodeGroup{
							&tree.Node{
								State: tree.NodeState{
									GameCount: 4,
									WinCount:  3,
								},
							},
							&tree.Node{
								State: tree.NodeState{
									GameCount: 6,
									WinCount:  5,
								},
							},
						}
					},
				},
				nodeSelector: MockNodeSelector{
					selectNode: func(
						nodes tree.NodeGroup,
					) *tree.Node {
						expectedNodes := tree.NodeGroup{
							&tree.Node{
								State: tree.NodeState{
									GameCount: 4,
									WinCount:  3,
								},
							},
							&tree.Node{
								State: tree.NodeState{
									GameCount: 6,
									WinCount:  5,
								},
							},
						}
						if !reflect.DeepEqual(
							nodes,
							expectedNodes,
						) {
							test.Fail()
						}

						return nodes[0]
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.NilPoint,
					},
					Board: func() models.Board {
						return models.NewBoard(
							models.Size{
								Width:  3,
								Height: 3,
							},
						)
					}(),
					State: tree.NodeState{
						GameCount: 2,
						WinCount:  1,
					},
				},
			},
			wantNode: &tree.Node{
				State: tree.NodeState{
					GameCount: 4,
					WinCount:  3,
				},
			},
			wantErr: nil,
		},
	} {
		searcher := MoveSearcher{
			MoveGenerator: data.fields.
				moveGenerator,
			Builder: data.fields.builder,
			NodeSelector: data.fields.
				nodeSelector,
		}
		gotNode, gotErr :=
			searcher.SearchMove(data.args.root)

		if !reflect.DeepEqual(
			gotNode,
			data.wantNode,
		) {
			test.Fail()
		}
		if gotErr != data.wantErr {
			test.Fail()
		}
	}
}
