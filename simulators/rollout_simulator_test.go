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
		moveSelector MoveSelector
	}
	type args struct {
		board models.Board
		color models.Color
	}
	type data struct {
		fields     fields
		args       args
		wantResult tree.GameResult
		wantCount  int
	}

	var iterationCount int
	for _, data := range []data{
		data{
			fields: fields{
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
						points := models.Size{
							Width:  3,
							Height: 3,
						}.Points()[iterationCount+1:]
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
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
				color: models.White,
			},
			wantResult: tree.Win,
			wantCount:  3,
		},
		data{
			fields: fields{
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
						points := models.Size{
							Width:  3,
							Height: 3,
						}.Points()[iterationCount:]
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
				board: func() models.Board {
					return models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)
				}(),
				color: models.White,
			},
			wantResult: tree.Loss,
			wantCount:  4,
		},
		data{
			fields: fields{
				moveSelector: nil,
			},
			args: args{
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
				color: models.White,
			},
			wantResult: tree.Loss,
			wantCount:  0,
		},
	} {
		iterationCount = 0

		simulator := RolloutSimulator{
			MoveSelector: data.fields.
				moveSelector,
		}
		gotResult := simulator.Simulate(
			data.args.board,
			data.args.color,
		)

		if gotResult != data.wantResult {
			test.Fail()
		}
		if iterationCount != data.wantCount {
			test.Fail()
		}
	}
}
