// +build long

package searchers_test

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
)

func TestSearch(test *testing.T) {
	type args struct {
		storage     models.StoneStorage
		color       models.Color
		maximalPass int
	}
	type data struct {
		args      args
		wantMoves []models.Move
		wantErr   error
	}

	settingGroup := []searchSettings{
		searchSettings{},
		searchSettings{
			parallelSimulator: true,
		},
		searchSettings{
			parallelBulkySimulator: true,
		},
		searchSettings{
			parallelBuilder: true,
		},
	}
	for _, settings := range settingGroup {
		for _, data := range []data{
			data{
				args: args{
					storage: func() models.StoneStorage {
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
					color:       models.Black,
					maximalPass: 2,
				},
				wantMoves: []models.Move{
					models.Move{},
				},
				wantErr: models.ErrAlreadyWin,
			},
			data{
				args: args{
					storage: func() models.StoneStorage {
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
					color:       models.Black,
					maximalPass: 1,
				},
				wantMoves: []models.Move{
					models.Move{},
				},
				wantErr: searchers.
					ErrFailedBuilding,
			},
			data{
				args: args{
					storage: func() models.StoneStorage {
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
					color:       models.Black,
					maximalPass: 2,
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
					storage: func() models.StoneStorage {
						board := models.NewBoard(
							models.Size{
								Width:  3,
								Height: 3,
							},
						)

						points := board.Size().Points()
						opponentIndex := len(points) - 3
						opponentPoint :=
							points[opponentIndex]
						board =
							board.ApplyMove(models.Move{
								Color: models.White,
								Point: opponentPoint,
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
					color:       models.Black,
					maximalPass: 1000,
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
			data{
				args: args{
					storage: func() models.StoneStorage {
						board := models.NewBoard(
							models.Size{
								Width:  3,
								Height: 4,
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
					color:       models.Black,
					maximalPass: 1000,
				},
				wantMoves: []models.Move{
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
							Column: 0,
							Row:    3,
						},
					},
					models.Move{
						Color: models.Black,
						Point: models.Point{
							Column: 2,
							Row:    3,
						},
					},
				},
				wantErr: nil,
			},
		} {
			settings.maximalPass =
				data.args.maximalPass

			gotMove, gotErr := search(
				data.args.storage,
				data.args.color,
				settings,
			)

			var hasMatch bool
			for _, move := range data.wantMoves {
				if reflect.DeepEqual(
					gotMove,
					move,
				) {
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
}
