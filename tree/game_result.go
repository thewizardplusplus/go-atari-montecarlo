package tree

// GameResult ...
type GameResult int

// ...
const (
	Loss GameResult = iota
	Win
)

// Invert ...
func (
	result GameResult,
) Invert() GameResult {
	if result == Loss {
		return Win
	}

	return Loss
}
