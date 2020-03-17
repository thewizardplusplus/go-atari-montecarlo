package tree

import (
	"math"
	"reflect"
	"testing"
)

func TestNewNodeState(test *testing.T) {
	type args struct {
		result GameResult
	}
	type data struct {
		args args
		want NodeState
	}

	for _, data := range []data{
		data{
			args: args{Loss},
			want: NodeState{
				GameCount: 1,
				WinCount:  0,
			},
		},
		data{
			args: args{Win},
			want: NodeState{
				GameCount: 1,
				WinCount:  1,
			},
		},
	} {
		got := NewNodeState(data.args.result)

		if !reflect.DeepEqual(
			got,
			data.want,
		) {
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
		data{
			fields: fields{
				gameCount: 4,
				winCount:  2,
			},
			want: 0.5,
		},
		data{
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

func TestNodeStateAddResult(
	test *testing.T,
) {
	type fields struct {
		gameCount int
		winCount  int
	}
	type args struct {
		result GameResult
	}
	type data struct {
		fields    fields
		args      args
		wantState NodeState
	}

	for _, data := range []data{
		data{
			fields: fields{
				gameCount: 4,
				winCount:  2,
			},
			args: args{Loss},
			wantState: NodeState{
				GameCount: 5,
				WinCount:  2,
			},
		},
		data{
			fields: fields{
				gameCount: 4,
				winCount:  2,
			},
			args: args{Win},
			wantState: NodeState{
				GameCount: 5,
				WinCount:  3,
			},
		},
	} {
		state := NodeState{
			GameCount: data.fields.gameCount,
			WinCount:  data.fields.winCount,
		}
		state.AddResult(data.args.result)

		if !reflect.DeepEqual(
			state,
			data.wantState,
		) {
			test.Fail()
		}
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
