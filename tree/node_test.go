package tree

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

type MockNodeSelector struct {
	selectNode func(nodes NodeGroup) *Node
}

func (selector MockNodeSelector) SelectNode(nodes NodeGroup) *Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}

func TestNodeShallowCopy(test *testing.T) {
	node := &Node{
		Parent: &Node{
			State: NodeState{
				GameCount: 2,
				WinCount:  1,
			},
		},
		Move: models.Move{
			Color: models.White,
			Point: models.Point{
				Column: 2,
				Row:    2,
			},
		},
		Storage: func() models.StoneStorage {
			board := models.NewBoard(
				models.Size{
					Width:  3,
					Height: 3,
				},
			)

			for _, move := range []models.Move{
				{
					Color: models.Black,
					Point: models.Point{
						Column: 0,
						Row:    0,
					},
				},
				{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
			} {
				board = board.ApplyMove(move)
			}

			return board
		}(),
		State: NodeState{
			GameCount: 4,
			WinCount:  3,
		},
		Children: NodeGroup{
			&Node{
				State: NodeState{
					GameCount: 6,
					WinCount:  5,
				},
			},
			&Node{
				State: NodeState{
					GameCount: 8,
					WinCount:  7,
				},
			},
		},
	}
	got := node.ShallowCopy()

	nodePointer := reflect.ValueOf(node).Pointer()
	gotPointer := reflect.ValueOf(got).Pointer()
	if gotPointer == nodePointer {
		test.Fail()
	}

	want := &Node{
		Move: models.Move{
			Color: models.White,
			Point: models.Point{
				Column: 2,
				Row:    2,
			},
		},
		Storage: func() models.StoneStorage {
			board := models.NewBoard(
				models.Size{
					Width:  3,
					Height: 3,
				},
			)

			for _, move := range []models.Move{
				{
					Color: models.Black,
					Point: models.Point{
						Column: 0,
						Row:    0,
					},
				},
				{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
			} {
				board = board.ApplyMove(move)
			}

			return board
		}(),
	}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}

func TestNodeUpdateState(test *testing.T) {
	type fields struct {
		parent *Node
		state  NodeState
	}
	type args struct {
		state NodeState
	}
	type data struct {
		fields   fields
		args     args
		wantNode *Node
	}

	for _, data := range []data{
		{
			fields: fields{
				parent: nil,
				state: NodeState{
					GameCount: 4,
					WinCount:  2,
				},
			},
			args: args{
				state: NodeState{
					GameCount: 3,
					WinCount:  2,
				},
			},
			wantNode: &Node{
				Parent: nil,
				State: NodeState{
					GameCount: 7,
					WinCount:  4,
				},
			},
		},
		{
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
			args: args{
				state: NodeState{
					GameCount: 3,
					WinCount:  2,
				},
			},
			wantNode: &Node{
				Parent: &Node{
					Parent: &Node{
						Parent: nil,
						State: NodeState{
							GameCount: 5,
							WinCount:  3,
						},
					},
					State: NodeState{
						GameCount: 6,
						WinCount:  3,
					},
				},
				State: NodeState{
					GameCount: 7,
					WinCount:  4,
				},
			},
		},
	} {
		node := &Node{
			Parent: data.fields.parent,
			State:  data.fields.state,
		}
		node.UpdateState(data.args.state)

		if !reflect.DeepEqual(node, data.wantNode) {
			test.Fail()
		}
	}
}

func TestNodeMergeChildren(test *testing.T) {
	type fields struct {
		state    NodeState
		children NodeGroup
	}
	type args struct {
		another *Node
	}
	type data struct {
		fields   fields
		args     args
		wantNode *Node
	}

	for _, data := range []data{
		{
			fields: fields{
				children: NodeGroup{
					&Node{
						Move: models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						State: NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					},
					&Node{
						Move: models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
						State: NodeState{
							GameCount: 4,
							WinCount:  3,
						},
					},
				},
			},
			args: args{
				another: &Node{
					Children: NodeGroup{
						&Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 0,
									Row:    0,
								},
							},
							State: NodeState{
								GameCount: 6,
								WinCount:  5,
							},
						},
						&Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 2,
									Row:    2,
								},
							},
							State: NodeState{
								GameCount: 8,
								WinCount:  7,
							},
						},
					},
				},
			},
			wantNode: &Node{
				Children: NodeGroup{
					&Node{
						Move: models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						State: NodeState{
							GameCount: 8,
							WinCount:  6,
						},
					},
					&Node{
						Move: models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
						State: NodeState{
							GameCount: 12,
							WinCount:  10,
						},
					},
				},
			},
		},
		{
			fields: fields{
				state: NodeState{
					GameCount: 2,
					WinCount:  1,
				},
				children: nil,
			},
			args: args{
				another: &Node{
					State: NodeState{
						GameCount: 10,
						WinCount:  2,
					},
					Children: NodeGroup{
						&Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 0,
									Row:    0,
								},
							},
							State: NodeState{
								GameCount: 4,
								WinCount:  3,
							},
						},
						&Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 2,
									Row:    2,
								},
							},
							State: NodeState{
								GameCount: 6,
								WinCount:  5,
							},
						},
					},
				},
			},
			wantNode: func() *Node {
				node := &Node{
					State: NodeState{
						GameCount: 2,
						WinCount:  1,
					},
				}
				node.Children = NodeGroup{
					&Node{
						Parent: node,
						Move: models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						State: NodeState{
							GameCount: 4,
							WinCount:  3,
						},
					},
					&Node{
						Parent: node,
						Move: models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
						State: NodeState{
							GameCount: 6,
							WinCount:  5,
						},
					},
				}

				return node
			}(),
		},
	} {
		node := &Node{
			State:    data.fields.state,
			Children: data.fields.children,
		}
		node.MergeChildren(data.args.another)

		if !reflect.DeepEqual(node, data.wantNode) {
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
		{
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
					selectNode: func(nodes NodeGroup) *Node {
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
		{
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
		got := node.SelectLeaf(data.args.selector)

		if !reflect.DeepEqual(got, data.want) {
			test.Fail()
		}
	}
}

func TestNodeExpandLeaf(test *testing.T) {
	type fields struct {
		move     models.Move
		storage  models.StoneStorage
		state    NodeState
		children NodeGroup
	}
	type args struct {
		generator models.Generator
	}
	type data struct {
		fields           fields
		args             args
		checkParents     bool
		wantOriginalNode *Node
		wantResultNodes  NodeGroup
	}

	for _, data := range []data{
		{
			fields: fields{
				move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					for _, move := range []models.Move{
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
					} {
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
			args: args{
				generator: models.MoveGenerator{},
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
				Storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					for _, move := range []models.Move{
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
					} {
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
			wantResultNodes: NodeGroup{
				&Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    2,
						},
					},
					Storage: func() models.StoneStorage {
						board := models.NewBoard(
							models.Size{
								Width:  3,
								Height: 3,
							},
						)

						for _, move := range []models.Move{
							{
								Color: models.Black,
								Point: models.Point{
									Column: 0,
									Row:    0,
								},
							},
							{
								Color: models.White,
								Point: models.Point{
									Column: 2,
									Row:    2,
								},
							},
						} {
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
		},
		{
			fields: fields{
				move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					for _, point := range board.Size().Points() {
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
			args: args{
				generator: models.MoveGenerator{},
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
				Storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					for _, point := range board.Size().Points() {
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
			wantResultNodes: NodeGroup{
				&Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    2,
						},
					},
					Storage: func() models.StoneStorage {
						board := models.NewBoard(
							models.Size{
								Width:  3,
								Height: 3,
							},
						)

						for _, point := range board.Size().Points() {
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
		},
		{
			fields: fields{
				move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					for _, move := range []models.Move{
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
					} {
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
			args: args{
				generator: models.MoveGenerator{},
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
				Storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					for _, move := range []models.Move{
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
					} {
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
					for _, move := range []models.Move{
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    0,
							},
						},
						{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    2,
							},
						},
					} {
						board = board.ApplyMove(move)
					}

					var children NodeGroup
					for _, move := range []models.Move{
						{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						{
							Color: models.Black,
							Point: models.Point{
								Column: 2,
								Row:    0,
							},
						},
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						{
							Color: models.Black,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    2,
							},
						},
						{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					} {
						childBoard := board.ApplyMove(move)
						child := &Node{
							Move:    move,
							Storage: childBoard,
						}
						children = append(children, child)
					}

					return children
				}(),
			},
			wantResultNodes: func() NodeGroup {
				board := models.NewBoard(
					models.Size{
						Width:  3,
						Height: 3,
					},
				)
				for _, move := range []models.Move{
					{
						Color: models.Black,
						Point: models.Point{
							Column: 0,
							Row:    0,
						},
					},
					{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    2,
						},
					},
				} {
					board = board.ApplyMove(move)
				}

				var children NodeGroup
				for _, move := range []models.Move{
					{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    0,
						},
					},
					{
						Color: models.Black,
						Point: models.Point{
							Column: 2,
							Row:    0,
						},
					},
					{
						Color: models.Black,
						Point: models.Point{
							Column: 0,
							Row:    1,
						},
					},
					{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    1,
						},
					},
					{
						Color: models.Black,
						Point: models.Point{
							Column: 2,
							Row:    1,
						},
					},
					{
						Color: models.Black,
						Point: models.Point{
							Column: 0,
							Row:    2,
						},
					},
					{
						Color: models.Black,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					},
				} {
					childBoard := board.ApplyMove(move)
					child := &Node{
						Move:    move,
						Storage: childBoard,
					}
					children = append(children, child)
				}

				return children
			}(),
		},
	} {
		node := &Node{
			Move:     data.fields.move,
			Storage:  data.fields.storage,
			State:    data.fields.state,
			Children: data.fields.children,
		}
		got := node.ExpandLeaf(data.args.generator)

		// check and reset parents
		if data.checkParents {
			for _, child := range node.Children {
				if !reflect.DeepEqual(child.Parent, node) {
					test.Fail()
				}

				child.Parent = nil
			}
		}

		if !reflect.DeepEqual(node, data.wantOriginalNode) {
			test.Fail()
		}
		if !reflect.DeepEqual(got, data.wantResultNodes) {
			test.Fail()
		}
	}
}
