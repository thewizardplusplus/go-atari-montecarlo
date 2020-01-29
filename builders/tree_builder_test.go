package builders

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

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

type MockSimulator struct {
	simulate func(
		board models.Board,
		color models.Color,
	) tree.GameResult
}

func (simulator MockSimulator) Simulate(
	board models.Board,
	color models.Color,
) tree.GameResult {
	if simulator.simulate == nil {
		panic("not implemented")
	}

	return simulator.simulate(board, color)
}

func TestTreeBuilderPass(test *testing.T) {
	type fields struct {
		nodeSelector tree.NodeSelector
		simulator    Simulator
	}
	type args struct {
		root *tree.Node
	}
	type data struct {
		fields   fields
		args     args
		wantRoot *tree.Node
	}

	for _, data := range []data{} {
		builder := TreeBuilder{
			NodeSelector: data.fields.nodeSelector,
			Simulator:    data.fields.simulator,
		}
		builder.Pass(data.args.root)

		if !reflect.DeepEqual(
			data.args.root,
			data.wantRoot,
		) {
			test.Fail()
		}
	}
}
