package simulators

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// MoveSelector ...
type MoveSelector interface {
	SelectMove(moves []models.Move) models.Move
}

// RolloutSimulator ...
type RolloutSimulator struct {
	MoveSelector MoveSelector
}

// Simulate ...
func (simulator RolloutSimulator) Simulate(
	board models.Board,
	color models.Color,
) tree.NodeState {
	startColor := color
	for {
		moves, err := board.LegalMoves(color)
		if err != nil {
			// no moves or an already finished game
			result := tree.NewGameResult(err)
			if color != startColor {
				result = result.Invert()
			}

			var state tree.NodeState
			state.AddResult(result)

			return state
		}

		move := simulator.MoveSelector.
			SelectMove(moves)

		board = board.ApplyMove(move)
		color = color.Negative()
	}
}
