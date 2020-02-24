package terminators

import (
	"testing"
)

func TestNewPassTerminator(
	test *testing.T,
) {
	terminator := NewPassTerminator(5)

	if terminator.maximalPass != 5 {
		test.Fail()
	}
}

func TestPassTerminatorIsSearchTerminated(
	test *testing.T,
) {
	type fields struct {
		maximalPass int
	}
	type args struct {
		pass int
	}
	type data struct {
		fields fields
		args   args
		want   bool
	}

	for _, data := range []data{
		data{
			fields: fields{5},
			args:   args{4},
			want:   false,
		},
		data{
			fields: fields{5},
			args:   args{5},
			want:   true,
		},
		data{
			fields: fields{5},
			args:   args{6},
			want:   true,
		},
	} {
		terminator := PassTerminator{
			maximalPass: data.fields.maximalPass,
		}
		got := terminator.IsSearchTerminated(
			data.args.pass,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
