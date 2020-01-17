package simulators

import (
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

// MoveSelector ...
type MoveSelector interface {
	SelectMove(moves []models.Move) models.Move
}

// Simulator ...
type Simulator struct {
	MoveSelector MoveSelector
}

// Simulate ...
func (simulator Simulator) Simulate(
	board models.Board,
	color models.Color,
) tree.GameResult {
	startColor := color
	for !board.HasCapture(color) {
		moves := board.Moves(color)
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
