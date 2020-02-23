package searchers_test

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
)

func TestSearch(test *testing.T) {
	type args struct {
		board     models.Board
		color     models.Color
		ucbFactor float64
		passCount int
	}
	type data struct {
		args      args
		wantMoves []models.Move
		wantErr   error
	}

	for _, data := range []data{
		data{
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
				color:     models.Black,
				ucbFactor: 1,
				passCount: 2,
			},
			wantMoves: []models.Move{
				models.Move{},
			},
			wantErr: models.ErrAlreadyLoss,
		},
		data{
			args: args{
				board: func() models.Board {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					points := board.Size().Points()
					points = points[:len(points)-1]

					for _, point := range points {
						move := models.Move{
							Color: models.White,
							Point: point,
						}
						board = board.ApplyMove(move)
					}

					return board
				}(),
				color:     models.Black,
				ucbFactor: 1,
				passCount: 1,
			},
			wantMoves: []models.Move{
				models.Move{},
			},
			wantErr: searchers.ErrFailedBuilding,
		},
		data{
			args: args{
				board: func() models.Board {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					points := board.Size().Points()
					points = points[:len(points)-1]

					for _, point := range points {
						move := models.Move{
							Color: models.White,
							Point: point,
						}
						board = board.ApplyMove(move)
					}

					return board
				}(),
				color:     models.Black,
				ucbFactor: 1,
				passCount: 2,
			},
			wantMoves: []models.Move{
				models.Move{
					Color: models.Black,
					Point: models.Point{
						Column: 2,
						Row:    2,
					},
				},
			},
			wantErr: nil,
		},
		data{
			args: args{
				board: func() models.Board {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					points := board.Size().Points()
					opponentIndex := len(points) - 3
					board =
						board.ApplyMove(models.Move{
							Color: models.White,
							Point: points[opponentIndex],
						})

					points = points[:opponentIndex]
					for _, point := range points {
						move := models.Move{
							Color: models.Black,
							Point: point,
						}
						board = board.ApplyMove(move)
					}

					return board
				}(),
				color:     models.Black,
				ucbFactor: 1,
				passCount: 1000,
			},
			wantMoves: []models.Move{
				models.Move{
					Color: models.Black,
					Point: models.Point{
						Column: 1,
						Row:    2,
					},
				},
			},
			wantErr: nil,
		},
	} {
		gotMove, gotErr := search(
			data.args.board,
			data.args.color,
			data.args.ucbFactor,
			data.args.passCount,
		)

		var hasMatch bool
		for _, move := range data.wantMoves {
			if reflect.DeepEqual(gotMove, move) {
				hasMatch = true
				break
			}
		}
		if !hasMatch {
			test.Fail()
		}

		if gotErr != data.wantErr {
			test.Fail()
		}
	}
}
