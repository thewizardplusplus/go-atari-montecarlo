package tree

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestNewNodeGroup(test *testing.T) {
	type args struct {
		moves   []models.Move
		options []NodeGroupOption
	}
	type data struct {
		args args
		want NodeGroup
	}

	for _, data := range []data{
		data{
			args: args{
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
				options: []NodeGroupOption{},
			},
			want: NodeGroup{
				&Node{
					Move: models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 2,
							Row:    0,
						},
					},
				},
				&Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    2,
						},
					},
				},
			},
		},
		data{
			args: args{
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
				options: []NodeGroupOption{
					WithParent(&Node{
						State: NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					}),
				},
			},
			want: NodeGroup{
				&Node{
					Parent: &Node{
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
				},
				&Node{
					Parent: &Node{
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
				},
			},
		},
		data{
			args: args{
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
				options: []NodeGroupOption{
					WithBoard(func() models.Board {
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
					}()),
				},
			},
			want: NodeGroup{
				&Node{
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
			data.args.moves,
			data.args.options...,
		)

		if !reflect.DeepEqual(
			got,
			data.want,
		) {
			test.Fail()
		}
	}
}

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
