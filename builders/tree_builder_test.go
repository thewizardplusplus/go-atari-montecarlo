package builders

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

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

func TestTreeBuilderPass(test *testing.T) {
	type fields struct {
		nodeSelector tree.NodeSelector
		simulator    simulators.Simulator
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields   fields
		args     args
		wantRoot *tree.Node
	}

	for _, data := range []data{
		data{
			fields: fields{
				nodeSelector: MockNodeSelector{
					selectNode: func(
						nodes tree.NodeGroup,
					) *tree.Node {
						return nodes[0]
					},
				},
				simulator: MockSimulator{
					simulate: func(
						root *tree.Node,
					) tree.NodeState {
						return tree.NodeState{
							GameCount: 1,
							WinCount:  1,
						}
					},
				},
			},
			args: args{
				root: func() *tree.Node {
					root := &tree.Node{
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

							moves := []models.Move{
								models.Move{
									Color: models.Black,
									Point: models.Point{
										Column: 0,
										Row:    0,
									},
								},
								models.Move{
									Color: models.White,
									Point: models.Point{
										Column: 2,
										Row:    2,
									},
								},
							}
							for _, move := range moves {
								board =
									board.ApplyMove(move)
							}

							return board
						}(),
						State: tree.NodeState{
							GameCount: 5,
							WinCount:  2,
						},
					}
					childOne := &tree.Node{
						Parent: root,
						Move: models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						Board: root.Board.
							ApplyMove(models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						}),
						State: tree.NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					}
					childTwo := &tree.Node{
						Parent: root,
						Move: models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
						Board: root.Board.
							ApplyMove(models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						}),
						State: tree.NodeState{
							GameCount: 3,
							WinCount:  2,
						},
					}
					root.Children = tree.NodeGroup{
						childOne,
						childTwo,
					}

					return root
				}(),
			},
			wantRoot: func() *tree.Node {
				root := &tree.Node{
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

						moves := []models.Move{
							models.Move{
								Color: models.Black,
								Point: models.Point{
									Column: 0,
									Row:    0,
								},
							},
							models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 2,
									Row:    2,
								},
							},
						}
						for _, move := range moves {
							board =
								board.ApplyMove(move)
						}

						return board
					}(),
					State: tree.NodeState{
						GameCount: 6,
						WinCount:  2,
					},
				}

				childOne := &tree.Node{
					Parent: root,
					Move: models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    0,
						},
					},
					Board: root.Board.
						ApplyMove(models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    0,
						},
					}),
					State: tree.NodeState{
						GameCount: 3,
						WinCount:  2,
					},
				}
				childTwo := &tree.Node{
					Parent: root,
					Move: models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					},
					Board: root.Board.
						ApplyMove(models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					}),
					State: tree.NodeState{
						GameCount: 3,
						WinCount:  2,
					},
				}

				childOneOne := &tree.Node{
					Parent: childOne,
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    0,
						},
					},
					Board: childOne.Board.
						ApplyMove(models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    0,
						},
					}),
					State: tree.NodeState{
						GameCount: 1,
					},
				}
				childOneTwo := &tree.Node{
					Parent: childOne,
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    1,
						},
					},
					Board: childOne.Board.
						ApplyMove(models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    1,
						},
					}),
				}
				childOneThree := &tree.Node{
					Parent: childOne,
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 1,
							Row:    1,
						},
					},
					Board: childOne.Board.
						ApplyMove(models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 1,
							Row:    1,
						},
					}),
				}
				childOneFour := &tree.Node{
					Parent: childOne,
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    1,
						},
					},
					Board: childOne.Board.
						ApplyMove(models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    1,
						},
					}),
				}
				childOneFive := &tree.Node{
					Parent: childOne,
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    2,
						},
					},
					Board: childOne.Board.
						ApplyMove(models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    2,
						},
					}),
				}
				childOneSix := &tree.Node{
					Parent: childOne,
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					},
					Board: childOne.Board.
						ApplyMove(models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					}),
				}

				root.Children = tree.NodeGroup{
					childOne,
					childTwo,
				}
				childOne.Children = tree.NodeGroup{
					childOneOne,
					childOneTwo,
					childOneThree,
					childOneFour,
					childOneFive,
					childOneSix,
				}

				return root
			}(),
		},
	} {
		builder := TreeBuilder{
			NodeSelector: data.fields.
				nodeSelector,
			Simulator: data.fields.simulator,
		}
		builder.Pass(data.args.root)

		if !reflect.DeepEqual(
			data.args.root,
			data.wantRoot,
		) {
			test.Fail()
		}
	}
}
