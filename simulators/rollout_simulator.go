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
	MoveGenerator models.Generator
	MoveSelector  MoveSelector
}

// Simulate ...
func (simulator RolloutSimulator) Simulate(
	root *tree.Node,
) tree.NodeState {
	storage, previousMove, startColor :=
		root.Storage, root.Move, root.Move.Color
	for {
		moves, err := simulator.MoveGenerator.
			LegalMoves(storage, previousMove)
		if err != nil {
			// no moves or an already finished game
			state := tree.NewNodeState(err)
			if previousMove.Color != startColor {
				state = state.Invert()
			}

			return state
		}

		move := simulator.MoveSelector.
			SelectMove(moves)
		storage, previousMove =
			storage.ApplyMove(move), move
	}
}
