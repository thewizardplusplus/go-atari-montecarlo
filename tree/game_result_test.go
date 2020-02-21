package tree

import (
	"errors"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestNewGameResult(test *testing.T) {
	type args struct {
		err error
	}
	type data struct {
		args       args
		wantResult GameResult
		wantPanic  bool
	}

	for _, data := range []data{
		data{
			args: args{
				err: models.ErrAlreadyLoss,
			},
			wantResult: Loss,
			wantPanic:  false,
		},
		data{
			args: args{
				err: models.ErrAlreadyWin,
			},
			wantResult: Win,
			wantPanic:  false,
		},
		data{
			args: args{
				err: errors.New("dummy"),
			},
			wantResult: 0,
			wantPanic:  true,
		},
	} {
		var gotResult GameResult
		var hasPanic bool
		func() {
			defer func() {
				if err := recover(); err != nil {
					hasPanic = true
				}
			}()

			gotResult =
				NewGameResult(data.args.err)
		}()

		if gotResult != data.wantResult {
			test.Fail()
		}
		if hasPanic != data.wantPanic {
			test.Fail()
		}
	}
}

func TestGameResultInvert(test *testing.T) {
	type data struct {
		result GameResult
		want   GameResult
	}

	for _, data := range []data{
		data{
			result: Loss,
			want:   Win,
		},
		data{
			result: Win,
			want:   Loss,
		},
	} {
		got := data.result.Invert()

		if got != data.want {
			test.Fail()
		}
	}
}
