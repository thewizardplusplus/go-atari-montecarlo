package atarimontecarlo

import (
	"reflect"
	"testing"
)

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
						WinCount:  3,
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
