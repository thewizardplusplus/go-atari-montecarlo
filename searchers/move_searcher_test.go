package searchers

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
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

type MockNodeSelector struct {
	selectNode func(
		nodes tree.NodeGroup,
	) *tree.Node
}

func (selector MockNodeSelector) SelectNode(
	nodes tree.NodeGroup,
) *tree.Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}

func TestMoveSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		builder      builders.Builder
		nodeSelector tree.NodeSelector
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields   fields
		args     args
		wantNode *tree.Node
		wantOk   bool
	}

	for _, data := range []data{
		data{
			fields: fields{
				builder: MockBuilder{
					pass: func(root *tree.Node) {
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

						// don't change the passed node
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
			wantNode: nil,
			wantOk:   false,
		},
		data{
			fields: fields{
				builder: MockBuilder{
					pass: func(root *tree.Node) {
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

						// add few children
						// to the passed node
						root.Children = tree.NodeGroup{
							&tree.Node{
								State: tree.NodeState{
									GameCount: 4,
									WinCount:  3,
								},
							},
							&tree.Node{
								State: tree.NodeState{
									GameCount: 6,
									WinCount:  5,
								},
							},
						}
					},
				},
				nodeSelector: MockNodeSelector{
					selectNode: func(
						nodes tree.NodeGroup,
					) *tree.Node {
						expectedNodes := tree.NodeGroup{
							&tree.Node{
								State: tree.NodeState{
									GameCount: 4,
									WinCount:  3,
								},
							},
							&tree.Node{
								State: tree.NodeState{
									GameCount: 6,
									WinCount:  5,
								},
							},
						}
						if !reflect.DeepEqual(
							nodes,
							expectedNodes,
						) {
							test.Fail()
						}

						return nodes[0]
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
			wantNode: &tree.Node{
				State: tree.NodeState{
					GameCount: 4,
					WinCount:  3,
				},
			},
			wantOk: true,
		},
	} {
		searcher := MoveSearcher{
			Builder: data.fields.builder,
			NodeSelector: data.fields.
				nodeSelector,
		}
		gotNode, gotOk :=
			searcher.SearchMove(data.args.root)

		if !reflect.DeepEqual(
			gotNode,
			data.wantNode,
		) {
			test.Fail()
		}
		if gotOk != data.wantOk {
			test.Fail()
		}
	}
}
