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
) tree.GameResult {
	startColor := color
	for !board.HasCapture(color) {
		moves := board.Moves(color)
		if len(moves) == 0 {
			break
		}

		move := simulator.MoveSelector.
			SelectMove(moves)

		board = board.ApplyMove(move)
		color = color.Negative()
	}

	var result tree.GameResult
	if color == startColor {
		result = tree.Loss
	} else {
		result = tree.Win
	}

	return result
}
