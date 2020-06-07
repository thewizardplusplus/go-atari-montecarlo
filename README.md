# go-atari-montecarlo

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-atari-montecarlo?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-atari-montecarlo)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-atari-montecarlo)](https://goreportcard.com/report/github.com/thewizardplusplus/go-atari-montecarlo)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-atari-montecarlo.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-atari-montecarlo)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-atari-montecarlo/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-atari-montecarlo)

The library that implements an [Atari Go](https://senseis.xmp.net/?AtariGo) engine based on the Monte Carlo tree search algorithm.

_**Disclaimer:** this library was written directly on an Android smartphone with the AnGoIde IDE._

## Installation

```
$ go get github.com/thewizardplusplus/go-atari-montecarlo
```

## Benchmarks

Without parallelism:

```
BenchmarkSearch_5Passes-8                             	     300	   6041138 ns/op
BenchmarkSearch_10Passes-8                            	     100	  11545956 ns/op
BenchmarkSearch_15Passes-8                            	     100	  16472837 ns/op
BenchmarkSearch_20Passes-8                            	     100	  21192349 ns/op
```

With the parallel simulator:

```
BenchmarkSearch_5PassesAndParallelSimulator-8         	     100	  11482508 ns/op
BenchmarkSearch_10PassesAndParallelSimulator-8        	      50	  24193690 ns/op
BenchmarkSearch_15PassesAndParallelSimulator-8        	      30	  35422514 ns/op
BenchmarkSearch_20PassesAndParallelSimulator-8        	      30	  46065076 ns/op
```

With the parallel bulky simulator:

```
BenchmarkSearch_5PassesAndParallelBulkySimulator-8    	      50	  25558541 ns/op
BenchmarkSearch_10PassesAndParallelBulkySimulator-8   	      20	  59455409 ns/op
BenchmarkSearch_15PassesAndParallelBulkySimulator-8   	      20	  85564306 ns/op
BenchmarkSearch_20PassesAndParallelBulkySimulator-8   	      10	 186877684 ns/op
```

With the parallel builder:

```
BenchmarkSearch_5PassesAndParallelBuilder-8           	     100	  15637022 ns/op
BenchmarkSearch_10PassesAndParallelBuilder-8          	      50	  32202269 ns/op
BenchmarkSearch_15PassesAndParallelBuilder-8          	      30	  52235995 ns/op
BenchmarkSearch_20PassesAndParallelBuilder-8          	      20	  71661751 ns/op
```

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
