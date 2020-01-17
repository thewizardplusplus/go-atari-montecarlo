package simulators

import (
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

type MockMoveSelector struct {
	selectMove func(
		moves []models.Move,
	) models.Move
}

func (selector MockMoveSelector) SelectMove(
	moves []models.Move,
) models.Move {
	if selector.selectMove == nil {
		panic("not implemented")
	}

	return selector.selectMove(moves)
}

func TestRolloutSimulatorSimulate(
	test *testing.T,
) {

}
