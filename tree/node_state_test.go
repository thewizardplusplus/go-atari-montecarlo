package tree

import (
	"errors"
	"math"
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestNewNodeState(test *testing.T) {
	type args struct {
		err error
	}
	type data struct {
		args      args
		wantState NodeState
		wantPanic bool
	}

	for _, data := range []data{
		{
			args: args{
				err: models.ErrAlreadyLoss,
			},
			wantState: NodeState{
				GameCount: 1,
				WinCount:  0,
			},
			wantPanic: false,
		},
		{
			args: args{
				err: models.ErrAlreadyWin,
			},
			wantState: NodeState{
				GameCount: 1,
				WinCount:  1,
			},
			wantPanic: false,
		},
		{
			args: args{
				err: errors.New("dummy"),
			},
			wantState: NodeState{},
			wantPanic: true,
		},
	} {
		var gotState NodeState
		var hasPanic bool
		func() {
			defer func() {
				if err := recover(); err != nil {
					hasPanic = true
				}
			}()

			gotState = NewNodeState(data.args.err)
		}()

		if gotState != data.wantState {
			test.Fail()
		}
		if hasPanic != data.wantPanic {
			test.Fail()
		}
	}
}

func TestNodeStateWinRate(test *testing.T) {
	type fields struct {
		gameCount int
		winCount  int
	}
	type data struct {
		fields fields
		want   float64
	}

	for _, data := range []data{
		{
			fields: fields{
				gameCount: 4,
				winCount:  2,
			},
			want: 0.5,
		},
		{
			fields: fields{
				gameCount: 0,
				winCount:  0,
			},
			want: math.Inf(+1),
		},
	} {
		state := NodeState{
			GameCount: data.fields.gameCount,
			WinCount:  data.fields.winCount,
		}
		got := state.WinRate()

		if got != data.want {
			test.Fail()
		}
	}
}

func TestNodeStateInvert(test *testing.T) {
	state := NodeState{
		GameCount: 5,
		WinCount:  2,
	}
	got := state.Invert()

	want := NodeState{
		GameCount: 5,
		WinCount:  3,
	}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}

func TestNodeStateUpdate(test *testing.T) {
	update := NodeState{
		GameCount: 2,
		WinCount:  1,
	}
	state := NodeState{
		GameCount: 3,
		WinCount:  2,
	}
	state.Update(update)

	want := NodeState{
		GameCount: 5,
		WinCount:  3,
	}
	if !reflect.DeepEqual(state, want) {
		test.Fail()
	}
}
