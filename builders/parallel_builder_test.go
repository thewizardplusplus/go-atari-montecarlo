package builders

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestParallelBuilderPass(
	test *testing.T,
) {
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
				board = board.ApplyMove(move)
			}

			return board
		}(),
		Children: tree.NodeGroup{
			&tree.Node{
				Move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 0,
						Row:    2,
					},
				},
				State: tree.NodeState{
					GameCount: 2,
					WinCount:  1,
				},
			},
			&tree.Node{
				Move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    0,
					},
				},
				State: tree.NodeState{
					GameCount: 4,
					WinCount:  3,
				},
			},
		},
	}
	innerBuilder := MockBuilder{
		pass: func(root *tree.Node) {
			expectedRoot := &tree.Node{
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
			}
			if !reflect.DeepEqual(
				root,
				expectedRoot,
			) {
				test.Fail()
			}

			root.Children = tree.NodeGroup{
				&tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    2,
						},
					},
					State: tree.NodeState{
						GameCount: 6,
						WinCount:  5,
					},
				},
				&tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    0,
						},
					},
					State: tree.NodeState{
						GameCount: 8,
						WinCount:  7,
					},
				},
			}
		},
	}
	builder := ParallelBuilder{
		Builder:     innerBuilder,
		Concurrency: 10,
	}
	builder.Pass(root)

	wantRoot := &tree.Node{
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
		Children: tree.NodeGroup{
			&tree.Node{
				Move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 0,
						Row:    2,
					},
				},
				State: tree.NodeState{
					GameCount: 62,
					WinCount:  51,
				},
			},
			&tree.Node{
				Move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 2,
						Row:    0,
					},
				},
				State: tree.NodeState{
					GameCount: 84,
					WinCount:  73,
				},
			},
		},
	}
	if !reflect.DeepEqual(
		root,
		wantRoot,
	) {
		test.Fail()
	}
}
