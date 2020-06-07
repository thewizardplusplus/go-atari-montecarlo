package selectors

import (
	"math/rand"
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestRandomMoveSelectorSelectMove(test *testing.T) {
	// make the random generator deterministic for test reproducibility
	rand.Seed(1)

	got := RandomMoveSelector{}.
		SelectMove(
			[]models.Move{
				{
					Color: models.White,
					Point: models.Point{
						Column: 1,
						Row:    2,
					},
				},
				{
					Color: models.White,
					Point: models.Point{
						Column: 3,
						Row:    4,
					},
				},
			},
		)

	want := models.Move{
		Color: models.White,
		Point: models.Point{
			Column: 3,
			Row:    4,
		},
	}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}
