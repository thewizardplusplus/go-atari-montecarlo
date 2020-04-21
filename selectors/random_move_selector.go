package selectors

import (
	"math/rand"

	models "github.com/thewizardplusplus/go-atari-models"
)

// RandomMoveSelector ...
type RandomMoveSelector struct{}

// SelectMove ...
func (
	selector RandomMoveSelector,
) SelectMove(
	moves []models.Move,
) models.Move {
	index := rand.Intn(len(moves))
	return moves[index]
}
