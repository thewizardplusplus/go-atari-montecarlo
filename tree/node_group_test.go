package tree

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestNewNodeGroup(test *testing.T) {
	type args struct {
		parent *Node
		moves  []models.Move
	}
	type data struct {
		args args
		want NodeGroup
	}

	for _, data := range []data{
		data{
			args: args{
				parent: &Node{
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
						WinCount:  1,
					},
				},
				moves: []models.Move{
					models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 2,
							Row:    0,
						},
					},
					models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    2,
						},
					},
				},
			},
			want: NodeGroup{
				&Node{
					Parent: &Node{
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
						State: NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					},
					Move: models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 2,
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
									Column: 2,
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
				},
				&Node{
					Parent: &Node{
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
						State: NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					},
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
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
									Column: 0,
									Row:    2,
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
				},
			},
		},
	} {
		got := NewNodeGroup(
			data.args.parent,
			data.args.moves,
		)

		if !reflect.DeepEqual(
			got,
			data.want,
		) {
			test.Fail()
		}
	}
}

func TestNodeGroupMerge(test *testing.T) {
	type args struct {
		another NodeGroup
	}
	type data struct {
		nodes     NodeGroup
		args      args
		wantNodes NodeGroup
	}

	for _, data := range []data{
		data{
			nodes: NodeGroup{
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
			args: args{
				another: NodeGroup{
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
			wantNodes: NodeGroup{
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
		data{
			nodes: NodeGroup{
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
			args: args{
				another: NodeGroup{
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
								Column: 1,
								Row:    1,
							},
						},
						State: NodeState{
							GameCount: 8,
							WinCount:  7,
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
							GameCount: 10,
							WinCount:  9,
						},
					},
				},
			},
			wantNodes: NodeGroup{
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
						GameCount: 14,
						WinCount:  12,
					},
				},
			},
		},
		data{
			nodes: NodeGroup{
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
			args: args{
				another: NodeGroup{
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
				},
			},
			wantNodes: NodeGroup{
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
						GameCount: 4,
						WinCount:  3,
					},
				},
			},
		},
		data{
			nodes: NodeGroup{
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
			args: args{
				another: nil,
			},
			wantNodes: NodeGroup{
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
	} {
		data.nodes.Merge(data.args.another)

		if !reflect.DeepEqual(
			data.nodes,
			data.wantNodes,
		) {
			test.Fail()
		}
	}
}
