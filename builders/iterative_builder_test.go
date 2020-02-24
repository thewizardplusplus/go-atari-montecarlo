package builders

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockBuilder struct {
	pass func(root *tree.Node)
}

func (builder MockBuilder) Pass(
	root *tree.Node,
) {
	if builder.pass == nil {
		panic("not implemented")
	}

	builder.pass(root)
}

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

func TestIterativeBuilderPass(
	test *testing.T,
) {
	type fields struct {
		builder    Builder
		terminator terminators.
				BuildingTerminator
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields        fields
		args          args
		wantPassCount int
	}

	var passCount int
	for _, data := range []data{
		data{
			fields: fields{
				builder: MockBuilder{
					pass: func(root *tree.Node) {
						defer func() { passCount++ }()

						expectedRoot := &tree.Node{
							State: tree.NodeState{
								GameCount: 2,
								WinCount:  1,
							},
						}
						if !reflect.DeepEqual(
							root,
							expectedRoot,
						) {
							test.Fail()
						}
					},
				},
				terminator: MockBuildingTerminator{
					isSearchTerminated: func(
						pass int,
					) bool {
						if pass != passCount {
							test.Fail()
						}

						return pass >= 0
					},
				},
			},
			args: args{
				root: &tree.Node{
					State: tree.NodeState{
						GameCount: 2,
						WinCount:  1,
					},
				},
			},
			wantPassCount: 0,
		},
		data{
			fields: fields{
				builder: MockBuilder{
					pass: func(root *tree.Node) {
						defer func() { passCount++ }()

						expectedRoot := &tree.Node{
							State: tree.NodeState{
								GameCount: 2,
								WinCount:  1,
							},
						}
						if !reflect.DeepEqual(
							root,
							expectedRoot,
						) {
							test.Fail()
						}
					},
				},
				terminator: MockBuildingTerminator{
					isSearchTerminated: func(
						pass int,
					) bool {
						if pass != passCount {
							test.Fail()
						}

						return pass >= 5
					},
				},
			},
			args: args{
				root: &tree.Node{
					State: tree.NodeState{
						GameCount: 2,
						WinCount:  1,
					},
				},
			},
			wantPassCount: 5,
		},
	} {
		passCount = 0

		builder := IterativeBuilder{
			Builder:    data.fields.builder,
			Terminator: data.fields.terminator,
		}
		builder.Pass(data.args.root)

		if passCount != data.wantPassCount {
			test.Fail()
		}
	}
}
