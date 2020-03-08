package searchers

import (
	"reflect"
	"testing"

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

	for _, data := range []data{} {
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
