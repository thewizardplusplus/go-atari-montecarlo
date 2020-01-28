package tree

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

type MockNodeSelector struct {
	selectNode func(nodes NodeGroup) *Node
}

func (selector MockNodeSelector) SelectNode(
	nodes NodeGroup,
) *Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}

func TestNodeAddResult(test *testing.T) {
	type fields struct {
		parent *Node
		state  NodeState
	}
	type args struct {
		result GameResult
	}
	type data struct {
		fields   fields
		args     args
		wantNode *Node
	}

	for _, data := range []data{
		data{
			fields: fields{
				parent: nil,
				state: NodeState{
					GameCount: 4,
					WinCount:  2,
				},
			},
			args: args{Win},
			wantNode: &Node{
				Parent: nil,
				State: NodeState{
					GameCount: 5,
					WinCount:  3,
				},
			},
		},
		data{
			fields: fields{
				parent: &Node{
					Parent: &Node{
						Parent: nil,
						State: NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					},
					State: NodeState{
						GameCount: 3,
						WinCount:  2,
					},
				},
				state: NodeState{
					GameCount: 4,
					WinCount:  2,
				},
			},
			args: args{Win},
			wantNode: &Node{
				Parent: &Node{
					Parent: &Node{
						Parent: nil,
						State: NodeState{
							GameCount: 3,
							WinCount:  2,
						},
					},
					State: NodeState{
						GameCount: 4,
						WinCount:  2,
					},
				},
				State: NodeState{
					GameCount: 5,
					WinCount:  3,
				},
			},
		},
	} {
		node := &Node{
			Parent: data.fields.parent,
			State:  data.fields.state,
		}
		node.AddResult(data.args.result)

		if !reflect.DeepEqual(
			node,
			data.wantNode,
		) {
			test.Fail()
		}
	}
}

func TestNodeSelectLeaf(test *testing.T) {
	type fields struct {
		state    NodeState
		children NodeGroup
	}
	type args struct {
		selector NodeSelector
	}
	type data struct {
		fields fields
		args   args
		want   *Node
	}

	for _, data := range []data{
		data{
			fields: fields{
				state: NodeState{
					GameCount: 1,
					WinCount:  2,
				},
				children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
						Children: NodeGroup{
							&Node{
								State: NodeState{
									GameCount: 5,
									WinCount:  6,
								},
							},
							&Node{
								State: NodeState{
									GameCount: 7,
									WinCount:  8,
								},
							},
						},
					},
					&Node{
						State: NodeState{
							GameCount: 9,
							WinCount:  10,
						},
					},
				},
			},
			args: args{
				selector: MockNodeSelector{
					selectNode: func(
						nodes NodeGroup,
					) *Node {
						checkOne := reflect.DeepEqual(
							nodes,
							NodeGroup{
								&Node{
									State: NodeState{
										GameCount: 3,
										WinCount:  4,
									},
									Children: NodeGroup{
										&Node{
											State: NodeState{
												GameCount: 5,
												WinCount:  6,
											},
										},
										&Node{
											State: NodeState{
												GameCount: 7,
												WinCount:  8,
											},
										},
									},
								},
								&Node{
									State: NodeState{
										GameCount: 9,
										WinCount:  10,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							nodes,
							NodeGroup{
								&Node{
									State: NodeState{
										GameCount: 5,
										WinCount:  6,
									},
								},
								&Node{
									State: NodeState{
										GameCount: 7,
										WinCount:  8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return nodes[0]
					},
				},
			},
			want: &Node{
				State: NodeState{
					GameCount: 5,
					WinCount:  6,
				},
			},
		},
		data{
			fields: fields{
				state: NodeState{
					GameCount: 1,
					WinCount:  2,
				},
				children: nil,
			},
			args: args{
				selector: nil,
			},
			want: &Node{
				State: NodeState{
					GameCount: 1,
					WinCount:  2,
				},
			},
		},
	} {
		node := &Node{
			State:    data.fields.state,
			Children: data.fields.children,
		}
		got :=
			node.SelectLeaf(data.args.selector)

		if !reflect.DeepEqual(
			got,
			data.want,
		) {
			test.Fail()
		}
	}
}

func TestNodeExpandLeaf(test *testing.T) {
	type fields struct {
		move     models.Move
		board    models.Board
		state    NodeState
		children NodeGroup
	}
	type data struct {
		fields           fields
		checkParents     bool
		wantOriginalNode *Node
		wantResultNode   *Node
	}

	for _, data := range []data{
		data{
			fields: fields{
				move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
				board: func() models.Board {
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				state: NodeState{
					GameCount: 0,
				},
				children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
			checkParents: false,
			wantOriginalNode: &Node{
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				State: NodeState{
					GameCount: 0,
				},
				Children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
			wantResultNode: &Node{
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				State: NodeState{
					GameCount: 0,
				},
				Children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
		},
		data{
			fields: fields{
				move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
				board: func() models.Board {
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
				state: NodeState{
					GameCount: 2,
				},
				children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
			checkParents: false,
			wantOriginalNode: &Node{
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
				State: NodeState{
					GameCount: 2,
				},
				Children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
			wantResultNode: &Node{
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
				State: NodeState{
					GameCount: 2,
				},
				Children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
		},
		data{
			fields: fields{
				move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
				board: func() models.Board {
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				state: NodeState{
					GameCount: 2,
				},
				children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 1,
							WinCount:  2,
						},
					},
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
					},
				},
			},
			checkParents: true,
			wantOriginalNode: &Node{
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				State: NodeState{
					GameCount: 2,
				},
				Children: func() NodeGroup {
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
						board = board.ApplyMove(move)
					}

					var children NodeGroup
					moves = []models.Move{
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 2,
								Row:    0,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    2,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board := board.ApplyMove(move)
						child := &Node{
							Move:  move,
							Board: board,
						}
						children = append(
							children,
							child,
						)
					}

					return children
				}(),
			},
			wantResultNode: &Node{
				Move: models.Move{
					Color: models.Black,
					Point: models.Point{
						Column: 1,
						Row:    0,
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
							Color: models.Black,
							Point: models.Point{
								Column: 1,
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				State: NodeState{
					GameCount: 0,
				},
				Children: nil,
			},
		},
	} {
		node := &Node{
			Move:     data.fields.move,
			Board:    data.fields.board,
			State:    data.fields.state,
			Children: data.fields.children,
		}
		got := node.ExpandLeaf()

		// check and reset parents
		if data.checkParents {
			for _, child := range node.Children {
				if !reflect.DeepEqual(
					child.Parent,
					node,
				) {
					test.Fail()
				}

				child.Parent = nil
			}
		}

		if !reflect.DeepEqual(
			node,
			data.wantOriginalNode,
		) {
			test.Fail()
		}
		if !reflect.DeepEqual(
			got,
			data.wantResultNode,
		) {
			test.Fail()
		}
	}
}
