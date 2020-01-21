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

	for _, data := range []data{
		data{
			fields: fields{
				state: NodeState{
					GameCount: 1,
					WinCount:  2,
				},
				children: NodeGroup{
					&Node{
						State: NodeState{
							GameCount: 3,
							WinCount:  4,
						},
						Children: NodeGroup{
							&Node{
								State: NodeState{
									GameCount: 5,
									WinCount:  6,
								},
							},
							&Node{
								State: NodeState{
									GameCount: 7,
									WinCount:  8,
								},
							},
						},
					},
					&Node{
						State: NodeState{
							GameCount: 9,
							WinCount:  10,
						},
					},
				},
			},
			args: args{
				selector: MockNodeSelector{
					selectNode: func(
						nodes NodeGroup,
					) *Node {
						checkOne := reflect.DeepEqual(
							nodes,
							NodeGroup{
								&Node{
									State: NodeState{
										GameCount: 3,
										WinCount:  4,
									},
									Children: NodeGroup{
										&Node{
											State: NodeState{
												GameCount: 5,
												WinCount:  6,
											},
										},
										&Node{
											State: NodeState{
												GameCount: 7,
												WinCount:  8,
											},
										},
									},
								},
								&Node{
									State: NodeState{
										GameCount: 9,
										WinCount:  10,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							nodes,
							NodeGroup{
								&Node{
									State: NodeState{
										GameCount: 5,
										WinCount:  6,
									},
								},
								&Node{
									State: NodeState{
										GameCount: 7,
										WinCount:  8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return nodes[0]
					},
				},
			},
			want: &Node{
				State: NodeState{
					GameCount: 5,
					WinCount:  6,
				},
			},
		},
		data{
			fields: fields{
				state: NodeState{
					GameCount: 1,
					WinCount:  2,
				},
				children: nil,
			},
			args: args{
				selector: nil,
			},
			want: &Node{
				State: NodeState{
					GameCount: 1,
					WinCount:  2,
				},
			},
		},
	} {
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
