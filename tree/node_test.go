package tree

import (
	"reflect"
	"testing"
)

type MockNodeSelector struct {
	selectNode func(nodes NodeGroup) *Node
}

func (selector MockNodeSelector) SelectNode(
	nodes NodeGroup,
) *Node {
	if selector.selectNode == nil {
		panic("not implemented")
	}

	return selector.selectNode(nodes)
}

func TestNodeSelectLeaf(test *testing.T) {
	type fields struct {
		state    NodeState
		children NodeGroup
	}
	type args struct {
		selector NodeSelector
	}
	type data struct {
		fields fields
		args   args
		want   *Node
	}

	for _, data := range []data{} {
		node := &Node{
			State:    data.fields.state,
			Children: data.fields.children,
		}
		got :=
			node.SelectLeaf(data.args.selector)

		if !reflect.DeepEqual(
			got,
			data.want,
		) {
			test.Fail()
		}
	}
}

func TestNodeAddResult(test *testing.T) {
	type fields struct {
		parent *Node
		state  NodeState
	}
	type args struct {
		result GameResult
	}
	type data struct {
		fields   fields
		args     args
		wantNode *Node
	}

	for _, data := range []data{
		data{
			fields: fields{
				parent: nil,
				state: NodeState{
					GameCount: 4,
					WinCount:  2,
				},
			},
			args: args{Win},
			wantNode: &Node{
				Parent: nil,
				State: NodeState{
					GameCount: 5,
					WinCount:  3,
				},
			},
		},
		data{
			fields: fields{
				parent: &Node{
					Parent: &Node{
						Parent: nil,
						State: NodeState{
							GameCount: 2,
							WinCount:  1,
						},
					},
					State: NodeState{
						GameCount: 3,
						WinCount:  2,
					},
				},
				state: NodeState{
					GameCount: 4,
					WinCount:  2,
				},
			},
			args: args{Win},
			wantNode: &Node{
				Parent: &Node{
					Parent: &Node{
						Parent: nil,
						State: NodeState{
							GameCount: 3,
							WinCount:  2,
						},
					},
					State: NodeState{
						GameCount: 4,
						WinCount:  2,
					},
				},
				State: NodeState{
					GameCount: 5,
					WinCount:  3,
				},
			},
		},
	} {
		node := &Node{
			Parent: data.fields.parent,
			State:  data.fields.state,
		}
		node.AddResult(data.args.result)

		if !reflect.DeepEqual(
			node,
			data.wantNode,
		) {
			test.Fail()
		}
	}
}
