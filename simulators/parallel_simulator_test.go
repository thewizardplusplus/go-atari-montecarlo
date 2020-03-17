package simulators

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

func TestParallelSimulatorSimulate(
	test *testing.T,
) {
	type fields struct {
		simulator   Simulator
		concurrency int
	}
	type args struct {
		board models.Board
		color models.Color
	}
	type data struct {
		fields    fields
		args      args
		wantState tree.NodeState
	}

	for _, data := range []data{} {
		simulator := ParallelSimulator{
			Simulator:   data.fields.simulator,
			Concurrency: data.fields.concurrency,
		}
		gotState := simulator.Simulate(
			data.args.board,
			data.args.color,
		)

		if !reflect.DeepEqual(
			gotState,
			data.wantState,
		) {
			test.Fail()
		}
	}
}
