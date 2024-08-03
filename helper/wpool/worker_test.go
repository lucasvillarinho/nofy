package wpool

import (
	"fmt"
	"testing"
)

func TestCollectResults(t *testing.T) {
	// Given a worker pool with some results
	wp := &WorkPool[int, int]{
		results: make(chan Task[int, int]),
	}
	t.Run("Collect results from workers", func(t *testing.T) {
		// given a worker pool with some results
		go func() {
			wp.results <- Task[int, int]{ID: "1", Input: 1, Result: 2}
			wp.results <- Task[int, int]{ID: "2", Input: 2, Result: 4}
			wp.results <- Task[int, int]{ID: "3", Input: 3, Result: 6}
			wp.results <- Task[int, int]{ID: "4", Input: 4, Result: 0, Err: fmt.Errorf("some error")}
			close(wp.results)
		}()
		// The expected results are the tasks that succeeded
		expected := map[string]Task[int, int]{
			"1": {ID: "1", Input: 1, Result: 2},
			"2": {ID: "2", Input: 2, Result: 4},
			"3": {ID: "3", Input: 3, Result: 6},
		}

		// when we collect the results
		results := collectResults(wp)

		// then we should get the expected results
		if len(results) != len(expected) {
			t.Errorf("expected %d results, got %d", len(expected), len(results))
		}

		// and the results should match the expected results
		for _, result := range results {
			expectedResult, exists := expected[result.ID]
			if !exists {
				t.Errorf("unexpected result %+v", result)
				continue
			}
			if result != expectedResult {
				t.Errorf("expected %+v, got %+v", expectedResult, result)
			}
		}
	})
}

func TestCollectErrors(t *testing.T) {
	// When a worker fails to process a task, it should send the task to the errors channel
	wp := &WorkPool[int, int]{
		errors: make(chan Task[int, int]),
	}

	t.Run("Collect errors from workers", func(t *testing.T) {
		// given a worker pool with some errors
		go func() {
			wp.errors <- Task[int, int]{ID: "1", Input: 1, Result: 0, Err: fmt.Errorf("some error 1")}
			wp.errors <- Task[int, int]{ID: "2", Input: 2, Result: 4}
			wp.errors <- Task[int, int]{ID: "3", Input: 3, Result: 0, Err: fmt.Errorf("some error 2")}
			wp.errors <- Task[int, int]{ID: "4", Input: 4, Result: 0, Err: fmt.Errorf("some error 3")}
			close(wp.errors)
		}()

		// The expected results are the tasks that failed
		expected := map[string]Task[int, int]{
			"1": {ID: "1", Input: 1, Result: 0, Err: fmt.Errorf("some error 1")},
			"3": {ID: "3", Input: 3, Result: 0, Err: fmt.Errorf("some error 2")},
			"4": {ID: "4", Input: 4, Result: 0, Err: fmt.Errorf("some error 3")},
		}

		// when we collect the errors
		errors := collectErrors(wp)

		// then we should get the expected results
		if len(errors) != len(expected) {
			t.Errorf("expected %d results, got %d", len(expected), len(errors))
		}

		// and the results should match the expected results
		for _, result := range errors {
			expectedResult, exists := expected[result.ID]
			if !exists {
				t.Errorf("unexpected result %+v", result)
				continue
			}
			if result.Err.Error() != expectedResult.Err.Error() {
				t.Errorf("expected %+v, got %+v", expectedResult, result)
			}
		}
	})
}
