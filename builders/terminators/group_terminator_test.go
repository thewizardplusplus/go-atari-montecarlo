package terminators

import (
	"reflect"
	"testing"
)

type MockBuildingTerminator struct {
	isSearchTerminated func(pass int) bool
}

func (
	terminator MockBuildingTerminator,
) IsSearchTerminated(pass int) bool {
	if terminator.isSearchTerminated == nil {
		panic("not implemented")
	}

	return terminator.isSearchTerminated(pass)
}

func TestNewGroupTerminator(
	test *testing.T,
) {
	type args struct {
		terminators []BuildingTerminator
	}
	type data struct {
		args args
	}

	for _, data := range []data{
		data{
			args: args{nil},
		},
		data{
			args: args{
				terminators: []BuildingTerminator{
					MockBuildingTerminator{
						isSearchTerminated: func(
							pass int,
						) bool {
							panic("not implemented")
						},
					},
					MockBuildingTerminator{
						isSearchTerminated: func(
							pass int,
						) bool {
							panic("not implemented")
						},
					},
				},
			},
		},
	} {
		group := NewGroupTerminator(
			data.args.terminators...,
		)

		if !reflect.DeepEqual(
			group.terminators,
			data.args.terminators,
		) {
			test.Fail()
		}
	}
}

func TestGroupTerminatorIsSearchTerminated(
	test *testing.T,
) {
	type fields struct {
		terminators []BuildingTerminator
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
			fields: fields{nil},
			args:   args{5},
			want:   false,
		},
		data{
			fields: fields{
				terminators: []BuildingTerminator{
					MockBuildingTerminator{
						isSearchTerminated: func(
							pass int,
						) bool {
							if pass != 5 {
								test.Fail()
							}

							return false
						},
					},
					MockBuildingTerminator{
						isSearchTerminated: func(
							pass int,
						) bool {
							if pass != 5 {
								test.Fail()
							}

							return false
						},
					},
				},
			},
			args: args{5},
			want: false,
		},
		data{
			fields: fields{
				terminators: []BuildingTerminator{
					MockBuildingTerminator{
						isSearchTerminated: func(
							pass int,
						) bool {
							if pass != 5 {
								test.Fail()
							}

							return true
						},
					},
					MockBuildingTerminator{
						isSearchTerminated: func(
							pass int,
						) bool {
							panic("not implemented")
						},
					},
				},
			},
			args: args{5},
			want: true,
		},
	} {
		group := GroupTerminator{
			terminators: data.fields.terminators,
		}
		got := group.IsSearchTerminated(
			data.args.pass,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
