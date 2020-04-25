package parallel

import (
	"sync"
)

// Task ...
type Task func(
	index int,
) (result interface{})

// Run ...
func Run(
	concurrency int,
	task Task,
) (results []interface{}) {
	var waiter sync.WaitGroup
	waiter.Add(concurrency)

	results = make([]interface{}, concurrency)
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer waiter.Done()
			results[index] = task(index)
		}(i)
	}

	waiter.Wait()
	return results
}
