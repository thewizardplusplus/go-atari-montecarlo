package searchers

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type MockSearcher struct {
	searchMove func(
		root *tree.Node,
	) (*tree.Node, error)
}

func (searcher MockSearcher) SearchMove(
	root *tree.Node,
) (*tree.Node, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(root)
}

func TestNewReusedSearcher(
	test *testing.T,
) {
	var innerSearcher MockSearcher
	searcher :=
		NewReusedSearcher(innerSearcher)

	if !reflect.DeepEqual(
		searcher.searcher,
		innerSearcher,
	) {
		test.Fail()
	}
	if searcher.previousMove != nil {
		test.Fail()
	}
}

func TestReusedSearcherSearchPreviousMove(
	test *testing.T,
) {
	type fields struct {
		previousMove *tree.Node
	}
	type args struct {
		sample *tree.Node
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
				previousMove: &tree.Node{
					Children: tree.NodeGroup{
						&tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 0,
									Row:    2,
								},
							},
							State: tree.NodeState{
								GameCount: 2,
								WinCount:  1,
							},
						},
						&tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 1,
									Row:    2,
								},
							},
							State: tree.NodeState{
								GameCount: 4,
								WinCount:  3,
							},
						},
						&tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 2,
									Row:    2,
								},
							},
							State: tree.NodeState{
								GameCount: 6,
								WinCount:  5,
							},
						},
					},
				},
			},
			args: args{
				sample: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 1,
							Row:    2,
						},
					},
				},
			},
			wantNode: &tree.Node{
				Move: models.Move{
					Color: models.White,
					Point: models.Point{
						Column: 1,
						Row:    2,
					},
				},
				State: tree.NodeState{
					GameCount: 4,
					WinCount:  3,
				},
			},
			wantOk: true,
		},
		data{
			fields: fields{
				previousMove: &tree.Node{
					Children: tree.NodeGroup{
						&tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 0,
									Row:    2,
								},
							},
							State: tree.NodeState{
								GameCount: 2,
								WinCount:  1,
							},
						},
						&tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 1,
									Row:    2,
								},
							},
							State: tree.NodeState{
								GameCount: 4,
								WinCount:  3,
							},
						},
						&tree.Node{
							Move: models.Move{
								Color: models.White,
								Point: models.Point{
									Column: 2,
									Row:    2,
								},
							},
							State: tree.NodeState{
								GameCount: 6,
								WinCount:  5,
							},
						},
					},
				},
			},
			args: args{
				sample: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    3,
						},
					},
				},
			},
			wantNode: nil,
			wantOk:   false,
		},
		data{
			fields: fields{
				previousMove: &tree.Node{
					Children: nil,
				},
			},
			args: args{
				sample: &tree.Node{
					Move: models.Move{
						Color: models.White,
						Point: models.Point{
							Column: 2,
							Row:    3,
						},
					},
				},
			},
			wantNode: nil,
			wantOk:   false,
		},
	} {
		searcher := ReusedSearcher{
			previousMove: data.fields.
				previousMove,
		}
		gotNode, gotOk :=
			searcher.searchPreviousMove(
				data.args.sample,
			)

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
