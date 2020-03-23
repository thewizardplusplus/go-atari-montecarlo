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
	root *tree.Node,
) tree.NodeState {
	board := root.Board
	nextColor := root.Move.Color.Negative()
	startColor := nextColor
	for {
		moves, err := board.LegalMoves(nextColor)
		if err != nil {
			// no moves or an already finished game
			result := tree.NewGameResult(err)
			if nextColor != startColor {
				result = result.Invert()
			}

			return tree.NewNodeState(result)
		}

		move := simulator.MoveSelector.
			SelectMove(moves)

		board = board.ApplyMove(move)
		nextColor = nextColor.Negative()
	}
}
