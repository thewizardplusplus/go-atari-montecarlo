package tree

import (
	"testing"
)

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
