package tree

import (
	"reflect"
	"testing"
)

func TestNodeStateWinRate(test *testing.T) {
	state := NodeState{
		GameCount: 10,
		WinCount:  2,
	}
	score := state.WinRate()

	if score != 0.2 {
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
