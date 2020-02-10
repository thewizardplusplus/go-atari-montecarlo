package builders

import (
	"reflect"
	"testing"

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

func TestIterativeBuilderPass(
	test *testing.T,
) {
	type fields struct {
		builder   Builder
		passCount int
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
				passCount: 0,
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
				passCount: 5,
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
			Builder:   data.fields.builder,
			PassCount: data.fields.passCount,
		}
		builder.Pass(data.args.root)

		if passCount != data.wantPassCount {
			test.Fail()
		}
	}
}
