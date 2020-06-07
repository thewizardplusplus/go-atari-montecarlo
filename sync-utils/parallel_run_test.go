package syncutils

import (
	"reflect"
	"testing"
)

func TestParallelRun(test *testing.T) {
	results := ParallelRun(10, func(index int) (result interface{}) {
		return index
	})

	expectedResults := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(results, expectedResults) {
		test.Fail()
	}
}
