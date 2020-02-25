package main

import (
	"errors"
	"fmt"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

type settings struct {
	// random, winRate or ucb
	selectorType string
	ucbFactor    float64
	maximalPass  int
}

func search(
	root *tree.Node,
	settings settings,
) (*tree.Node, error) {
	var randomSelector selectors.RandomSelector
	var generalSelector tree.NodeSelector
	switch settings.selectorType {
	case "random":
		generalSelector = randomSelector
	case "winRate":
		generalSelector =
			selectors.MaximalSelector{
				NodeScorer: scorers.WinRateScorer{},
			}
	case "ucb":
		generalSelector =
			selectors.MaximalSelector{
				NodeScorer: scorers.UCBScorer{
					Factor: settings.ucbFactor,
				},
			}
	default:
		return nil,
			errors.New("unknown selector type")
	}

	simulator := simulators.RolloutSimulator{
		MoveSelector: selectors.MoveSelector{
			NodeSelector: randomSelector,
		},
	}
	terminator :=
		terminators.NewPassTerminator(
			settings.maximalPass,
		)
	builder := builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector: generalSelector,
			Simulator:    simulator,
		},
		Terminator: terminator,
	}
	searcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: generalSelector,
	}
	return searcher.SearchMove(root)
}

func main() {
	settings := settings{
		selectorType: "ucb",
		ucbFactor:    1,
		maximalPass:  100,
	}

	board := models.NewBoard(
		models.Size{
			Width:  5,
			Height: 5,
		},
	)
	root := tree.NewNode(board, models.Black)
	node, err := search(root, settings)

	fmt.Println(node, err)
}
