package simulators

import (
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
				moveSelector: MockMoveSelector{
					selectMove: func(
						moves []models.Move,
					) models.Move {
						defer func() {
							iterationCount++
						}()

						return moves[0]
					},
				},
			},
			args: args{
				board: func() models.Board {
					board := models.NewBoard(
						models.Size{
							Width:  5,
							Height: 5,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 2,
								Row:    3,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 3,
								Row:    2,
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
			wantCount:  5,
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
