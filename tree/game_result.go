package tree

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// GameResult ...
type GameResult int

// ...
const (
	Loss GameResult = iota
	Win
)

// NewGameResult ...
//
// Passed error should be
// models.ErrAlreadyLoss or
// models.ErrAlreadyWin only.
//
// Otherwize the function will panic.
func NewGameResult(err error) GameResult {
	var result GameResult
	switch err {
	case models.ErrAlreadyLoss:
		result = Loss
	case models.ErrAlreadyWin:
		result = Win
	default:
		panic(
			"tree.NewGameResult: " +
				"unsupported error",
		)
	}

	return result
}

// Invert ...
func (
	result GameResult,
) Invert() GameResult {
	if result == Loss {
		return Win
	}

	return Loss
}
