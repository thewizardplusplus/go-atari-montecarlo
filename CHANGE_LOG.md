# Change Log

## [v1.4](https://github.com/thewizardplusplus/go-atari-montecarlo/tree/v1.4) (2020-06-27)

- fixing the code style;
- adding usage examples;
- improving repository decor.

## [v1.3](https://github.com/thewizardplusplus/go-atari-montecarlo/tree/v1.3) (2020-05-29)

- use:
  - board interface;
  - separate move generator.

## [v1.2](https://github.com/thewizardplusplus/go-atari-montecarlo/tree/v1.2) (2020-05-03)

- remove:
  - separate representation of a game result;
  - node scoring by node win rate;
  - resetting of iteration terminators;
  - move searchers:
    - searcher that reuses a built tree;
    - fallback searcher that uses an additional searcher when the primary one returns an error;
- improve:
  - move selectors:
    - simplify:
      - constructing of node group;
      - architecture of move selectors;
    - optimize node scoring by the [Upper Confidence Bound algorithm](https://en.wikipedia.org/wiki/Multi-armed_bandit);
  - add a universal utility for parallel processing.

## [v1.1](https://github.com/thewizardplusplus/go-atari-montecarlo/tree/v1.1) (2020-03-30)

- optimization via parallel move searching:
  - parallel game simulating:
    - of a single node child;
    - of all node children;
  - parallel tree building;
- easily extensible and composable architecture:
  - all parallel algorithms can be combined in any combination.

## [v1.0](https://github.com/thewizardplusplus/go-atari-montecarlo/tree/v1.0) (2020-03-11)
