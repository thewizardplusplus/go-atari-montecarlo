package simulators

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockMoveSelector struct {
	selectMove func(
		moves []models.Move,
	) models.Move
}

func (selector MockMoveSelector) SelectMove(
	moves []models.Move,
) models.Move {
	if selector.selectMove == nil {
		panic("not implemented")
	}

	return selector.selectMove(moves)
}

func TestRolloutSimulatorSimulate(
	test *testing.T,
) {
	type fields struct {
		moveGenerator models.Generator
		moveSelector  MoveSelector
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields    fields
		args      args
		wantState tree.NodeState
		wantCount int
	}

	var iterationCount int
	for _, data := range []data{
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				// +--+--+--+
				// |B0|W1|B2|
				// +--+--+--+
				// |W3|  |  |
				// +--+--+--+
				// |  |  |  |
				// +--+--+--+
				moveSelector: MockMoveSelector{
					selectMove: func(
						moves []models.Move,
					) models.Move {
						defer func() {
							iterationCount++
						}()

						var color models.Color
						if iterationCount%2 == 0 {
							color = models.White
						} else {
							color = models.Black
						}

						var expectedMoves []models.Move
						size := models.Size{
							Width:  3,
							Height: 3,
						}
						points := size.Points()
						points =
							points[iterationCount+1:]
						for _, point := range points {
							move := models.Move{
								Color: color,
								Point: point,
							}
							expectedMoves =
								append(expectedMoves, move)
						}
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return moves[0]
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 0,
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
						}
						for _, move := range moves {
							board = board.ApplyMove(move)
						}

						return board
					}(),
				},
			},
			wantState: tree.NodeState{
				GameCount: 1,
				WinCount:  1,
			},
			wantCount: 3,
		},
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				// +--+--+--+
				// |W1|B2|W3|
				// +--+--+--+
				// |B4|  |  |
				// +--+--+--+
				// |  |  |  |
				// +--+--+--+
				moveSelector: MockMoveSelector{
					selectMove: func(
						moves []models.Move,
					) models.Move {
						defer func() {
							iterationCount++
						}()

						var color models.Color
						if iterationCount%2 == 0 {
							color = models.White
						} else {
							color = models.Black
						}

						var expectedMoves []models.Move
						size := models.Size{
							Width:  3,
							Height: 3,
						}
						points := size.Points()
						points = points[iterationCount:]
						for _, point := range points {
							move := models.Move{
								Color: color,
								Point: point,
							}
							expectedMoves =
								append(expectedMoves, move)
						}
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return moves[0]
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.Black,
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
				},
			},
			wantState: tree.NodeState{
				GameCount: 1,
				WinCount:  0,
			},
			wantCount: 4,
		},
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				moveSelector: MockMoveSelector{
					selectMove: func(
						moves []models.Move,
					) models.Move {
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
				},
			},
			wantState: tree.NodeState{
				GameCount: 1,
				WinCount:  1,
			},
			wantCount: 0,
		},
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				// +--+--+--+
				// |B0|W0|  |
				// +--+--+--+
				// |W0|  |  |
				// +--+--+--+
				// |  |  |  |
				// +--+--+--+
				moveSelector: MockMoveSelector{
					selectMove: func(
						moves []models.Move,
					) models.Move {
						panic("not implemented")
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 0,
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
								Color: models.White,
								Point: models.Point{
									Column: 1,
									Row:    0,
								},
							},
							models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 0,
									Row:    1,
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
			wantState: tree.NodeState{
				GameCount: 1,
				WinCount:  1,
			},
			wantCount: 0,
		},
		data{
			fields: fields{
				moveGenerator: models.
					MoveGenerator{},
				// +--+--+--+
				// |B0|W0|  |
				// +--+--+--+
				// |W0|  |  |
				// +--+--+--+
				// |  |  |  |
				// +--+--+--+
				moveSelector: MockMoveSelector{
					selectMove: func(
						moves []models.Move,
					) models.Move {
						panic("not implemented")
					},
				},
			},
			args: args{
				root: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 0,
							Row:    1,
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
									Column: 1,
									Row:    0,
								},
							},
							models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 0,
									Row:    1,
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
			wantState: tree.NodeState{
				GameCount: 1,
				WinCount:  0,
			},
			wantCount: 0,
		},
	} {
		iterationCount = 0

		simulator := RolloutSimulator{
			MoveGenerator: data.fields.
				moveGenerator,
			MoveSelector: data.fields.
				moveSelector,
		}
		gotState :=
			simulator.Simulate(data.args.root)

		if !reflect.DeepEqual(
			gotState,
			data.wantState,
		) {
			test.Fail()
		}
		if iterationCount != data.wantCount {
			test.Fail()
		}
	}
}
