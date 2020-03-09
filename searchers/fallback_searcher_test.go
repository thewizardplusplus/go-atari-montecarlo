package searchers

import (
	"errors"
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestFallbackSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		primarySearcher  Searcher
		fallbackSearcher Searcher
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields   fields
		args     args
		wantNode *tree.Node
		wantErr  error
	}

	for _, data := range []data{
		data{
			fields: fields{
				primarySearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						node := &tree.Node{
							State: tree.NodeState{
								GameCount: 4,
								WinCount:  3,
							},
						}
						return node, nil
					},
				},
				fallbackSearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
						panic("not implemented")
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
			wantErr: nil,
		},
		data{
			fields: fields{
				primarySearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						return nil,
							models.ErrAlreadyLoss
					},
				},
				fallbackSearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
						panic("not implemented")
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
			wantErr:  models.ErrAlreadyLoss,
		},
		data{
			fields: fields{
				primarySearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						return nil, models.ErrAlreadyWin
					},
				},
				fallbackSearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
						panic("not implemented")
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
			wantErr:  models.ErrAlreadyWin,
		},
		data{
			fields: fields{
				primarySearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						return nil, errors.New("dummy")
					},
				},
				fallbackSearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						node := &tree.Node{
							State: tree.NodeState{
								GameCount: 6,
								WinCount:  5,
							},
						}
						return node, nil
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
					GameCount: 6,
					WinCount:  5,
				},
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				primarySearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						return nil,
							errors.New("dummy #1")
					},
				},
				fallbackSearcher: MockSearcher{
					searchMove: func(
						root *tree.Node,
					) (*tree.Node, error) {
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

						return nil,
							errors.New("dummy #2")
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
			wantErr:  errors.New("dummy #2"),
		},
	} {
		searcher := FallbackSearcher{
			PrimarySearcher: data.fields.
				primarySearcher,
			FallbackSearcher: data.fields.
				fallbackSearcher,
		}
		gotNode, gotErr :=
			searcher.SearchMove(data.args.root)

		if !reflect.DeepEqual(
			gotNode,
			data.wantNode,
		) {
			test.Fail()
		}
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Fail()
		}
	}
}
