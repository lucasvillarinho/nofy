package wpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorkerStart(t *testing.T) {
	var wg sync.WaitGroup
	processFunc := func(input int) (int, error) {
		return input * 2, nil
	}
	errorFunc := func(input int) (int, error) {
		return 0, fmt.Errorf("processing error for input %d", input)
	}

	t.Run("Worker processed result successfully", func(t *testing.T) {
		jobCh := make(chan *Job[int, int], 1)
		quitCh := make(chan struct{})
		resultsCh := make(chan *Job[int, int], 1)
		errorsCh := make(chan *Job[int, int], 1)

		worker, err := NewWorker[int, int](
			WithJobChannel[int, int](jobCh),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](quitCh),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](resultsCh),
			WithErrorsChannel[int, int](errorsCh),
		)
		if err != nil {
			t.Fatalf("Failed to create worker: %v", err)
		}

		worker.start()

		testJob := &Job[int, int]{ID: "test-1", Input: 5}
		jobCh <- testJob

		time.Sleep(100 * time.Millisecond)

		close(quitCh)
		wg.Wait()

		if testJob.Result != 10 {
			t.Errorf("Expected result 10, but got %d", testJob.Result)
		}

		if testJob.Err != nil {
			t.Errorf("Expected error nil, but got %v", testJob.Err)
		}
	})
	t.Run("Worker processed result with error", func(t *testing.T) {
		jobCh := make(chan *Job[int, int], 1)
		quitCh := make(chan struct{})
		resultsCh := make(chan *Job[int, int], 1)
		errorsCh := make(chan *Job[int, int], 1)

		worker, err := NewWorker(
			WithJobChannel[int, int](jobCh),
			WithProcessFunc[int, int](errorFunc),
			WithQuitChannel[int, int](quitCh),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](resultsCh),
			WithErrorsChannel[int, int](errorsCh),
		)
		if err != nil {
			t.Fatalf("Failed to create worker: %v", err)
		}

		worker.start()

		testJob := &Job[int, int]{ID: "test-2", Input: 5}
		jobCh <- testJob

		time.Sleep(100 * time.Millisecond)

		close(quitCh)
		wg.Wait()

		if testJob.Result != 0 {
			t.Errorf("Expected result 0, but got %d", testJob.Result)
		}

		if testJob.Err == nil {
			t.Errorf("Expected err not nil, but got nil")
		}

		if testJob.Err.Error() != "processing error for input 5" {
			t.Errorf("Expected error message 'processing error for input 5', but got '%s'", testJob.Err.Error())
		}

	})
}


func TestNewWorker(t *testing.T) {
	var wg sync.WaitGroup

	processFunc := func(input int) (int, error) {
		return input * 2, nil
	}

	t.Run("Worker created successfully", func(t *testing.T) {
		jobCh := make(chan *Job[int, int], 1)
		quitCh := make(chan struct{})
		resultsCh := make(chan *Job[int, int], 1)
		errorsCh := make(chan *Job[int, int], 1)

		worker, err := NewWorker(
			WithJobChannel[int, int](jobCh),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](quitCh),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](resultsCh),
			WithErrorsChannel[int, int](errorsCh),
		)
		if err != nil {
			t.Fatalf("Failed to create worker: %v", err)
		}

		if worker == nil {
			t.Fatal("Expected worker to be created, but got nil")
		}
	})

	t.Run("Worker creation fails with nil job channel", func(t *testing.T) {
		_, err := NewWorker[int, int](
			WithJobChannel[int, int](nil),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](make(chan struct{})),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](make(chan *Job[int, int], 1)),
			WithErrorsChannel[int, int](make(chan *Job[int, int], 1)),
		)
		if err == nil {
			t.Fatal("Expected error when job channel is nil, but got nil")
		}
	})

	t.Run("Worker creation fails with nil process function", func(t *testing.T) {
		_, err := NewWorker[int, int](
			WithJobChannel[int, int](make(chan *Job[int, int], 1)),
			WithProcessFunc[int, int](nil),
			WithQuitChannel[int, int](make(chan struct{})),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](make(chan *Job[int, int], 1)),
			WithErrorsChannel[int, int](make(chan *Job[int, int], 1)),
		)
		if err == nil {
			t.Fatal("Expected error when process function is nil, but got nil")
		}
	})

	t.Run("Worker creation fails with nil quit channel", func(t *testing.T) {
		_, err := NewWorker[int, int](
			WithJobChannel[int, int](make(chan *Job[int, int], 1)),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](nil),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](make(chan *Job[int, int], 1)),
			WithErrorsChannel[int, int](make(chan *Job[int, int], 1)),
		)
		if err == nil {
			t.Fatal("Expected error when quit channel is nil, but got nil")
		}
	})

	t.Run("Worker creation fails with nil wait group", func(t *testing.T) {
		_, err := NewWorker[int, int](
			WithJobChannel[int, int](make(chan *Job[int, int], 1)),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](make(chan struct{})),
			WithWaitGroup[int, int](nil),
			WithResultsChannel[int, int](make(chan *Job[int, int], 1)),
			WithErrorsChannel[int, int](make(chan *Job[int, int], 1)),
		)
		if err == nil {
			t.Fatal("Expected error when wait group is nil, but got nil")
		}
	})

	t.Run("Worker creation fails with nil results channel", func(t *testing.T) {
		_, err := NewWorker[int, int](
			WithJobChannel[int, int](make(chan *Job[int, int], 1)),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](make(chan struct{})),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](nil),
			WithErrorsChannel[int, int](make(chan *Job[int, int], 1)),
		)
		if err == nil {
			t.Fatal("Expected error when results channel is nil, but got nil")
		}
	})

	t.Run("Worker creation fails with nil errors channel", func(t *testing.T) {
		_, err := NewWorker[int, int](
			WithJobChannel[int, int](make(chan *Job[int, int], 1)),
			WithProcessFunc[int, int](processFunc),
			WithQuitChannel[int, int](make(chan struct{})),
			WithWaitGroup[int, int](&wg),
			WithResultsChannel[int, int](make(chan *Job[int, int], 1)),
			WithErrorsChannel[int, int](nil),
		)
		if err == nil {
			t.Fatal("Expected error when errors channel is nil, but got nil")
		}
	})
}

func TestWithProcessFunc(t *testing.T) {
	t.Run("Returns error when process function is nil", func(t *testing.T) {
		opt := WithProcessFunc[int, int](nil)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
		expectedErr := "process function cannot be nil"
		if err.Error() != expectedErr {
			t.Fatalf("Expected error '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("Sets process function when it is not nil", func(t *testing.T) {
		processFunc := func(input int) (int, error) {
			return input * 2, nil
		}
		opt := WithProcessFunc(processFunc)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err != nil {
			t.Fatalf("Expected no error, but got '%s'", err.Error())
		}
		if worker.process == nil {
			t.Fatal("Expected process function to be set, but it is nil")
		}
	})
}



func TestWithQuitChannel(t *testing.T) {
	t.Run("Returns error when quit channel is nil", func(t *testing.T) {
		opt := WithQuitChannel[int, int](nil)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
		expectedErr := "quit channel cannot be nil"
		if err.Error() != expectedErr {
			t.Fatalf("Expected error '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("Sets quit channel when it is not nil", func(t *testing.T) {
		quitCh := make(chan struct{})
		opt := WithQuitChannel[int, int](quitCh)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err != nil {
			t.Fatalf("Expected no error, but got '%s'", err.Error())
		}
		if worker.quit == nil {
			t.Fatal("Expected quit channel to be set, but it is nil")
		}
	})
}



func TestWithResultsChannel(t *testing.T) {
	t.Run("Returns error when results channel is nil", func(t *testing.T) {
		opt := WithResultsChannel[int, int](nil)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
		expectedErr := "results channel cannot be nil"
		if err.Error() != expectedErr {
			t.Fatalf("Expected error '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("Sets results channel when it is not nil", func(t *testing.T) {
		resultsCh := make(chan *Job[int, int])
		opt := WithResultsChannel(resultsCh)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err != nil {
			t.Fatalf("Expected no error, but got '%s'", err.Error())
		}
		if worker.results == nil {
			t.Fatal("Expected results channel to be set, but it is nil")
		}
	})
}


func TestWithErrorsChannel(t *testing.T) {
	t.Run("Returns error when errors channel is nil", func(t *testing.T) {
		opt := WithErrorsChannel[int, int](nil)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
		expectedErr := "errors channel cannot be nil"
		if err.Error() != expectedErr {
			t.Fatalf("Expected error '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("Sets errors channel when it is not nil", func(t *testing.T) {
		errorsCh := make(chan *Job[int, int])
		opt := WithErrorsChannel(errorsCh)
		worker := &Worker[int, int]{}
		err := opt(worker)
		if err != nil {
			t.Fatalf("Expected no error, but got '%s'", err.Error())
		}
		if worker.errors == nil {
			t.Fatal("Expected errors channel to be set, but it is nil")
		}
	})
}