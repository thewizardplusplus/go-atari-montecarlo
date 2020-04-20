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
	board, previousMove, startColor :=
		root.Board, root.Move, root.Move.Color
	for {
		moves, err :=
			board.LegalMoves(previousMove)
		if err != nil {
			// no moves or an already finished game
			result := tree.NewGameResult(err)
			if previousMove.Color != startColor {
				result = result.Invert()
			}

			return tree.NewNodeState(result)
		}

		move := simulator.MoveSelector.
			SelectMove(moves)
		board, previousMove =
			board.ApplyMove(move), move
	}
}
