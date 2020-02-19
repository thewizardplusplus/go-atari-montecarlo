package searchers_test

import (
	"log"
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestSearch(test *testing.T) {
	test.Skip()

	type args struct {
		board     models.Board
		color     models.Color
		ucbFactor float64
		passCount int
	}
	type data struct {
		args     args
		wantMove models.Move
		wantOk   bool
	}

	for index, data := range []data{
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
				ucbFactor: 2,
				passCount: 2,
			},
			wantMove: models.Move{},
			wantOk:   false,
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
				ucbFactor: 2,
				passCount: 1,
			},
			wantMove: models.Move{},
			wantOk:   false,
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
				ucbFactor: 2,
				passCount: 2,
			},
			wantMove: models.Move{
				Color: models.Black,
				Point: models.Point{
					Column: 2,
					Row:    2,
				},
			},
			wantOk: true,
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

					moves := []models.Move{
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
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
							Color: models.White,
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
						board = board.ApplyMove(move)
					}

					return board
				}(),
				color:     models.Black,
				ucbFactor: 0.001,
				passCount: 1000,
			},
			wantMove: models.Move{
				Color: models.Black,
				Point: models.Point{
					Column: 1,
					Row:    0,
				},
			},
			wantOk: true,
		},
	} {
		log.Printf("run #%d", index)

		gotMove, gotOk := search(
			data.args.board,
			data.args.color,
			data.args.ucbFactor,
			data.args.passCount,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Log(gotMove)
			test.Log(data.wantMove)
			test.Fail()
		}
		if gotOk != data.wantOk {
			test.Log(gotOk)
			test.Log(data.wantOk)
			test.Fail()
		}
	}
}
