package wpool

import (
	"fmt"
	"sync"
)

// Task represents a task to be executed by a worker
type Task[T any, R any] struct {
	ID     string
	Input  T
	Result R
	Err    error
}

// WorkPool Worker represents a worker that executes tasks
type WorkPool[T any, R any] struct {
	numWorkers int
	jobs       chan Task[T, R]
	results    chan Task[T, R]
	errors     chan Task[T, R]
	process    func(T) (R, error)
	wg         sync.WaitGroup
	quit       chan struct{}
	once       sync.Once
}

// NewWorkPool creates a new worker pool
func NewWorkPool[T any, R any](numWorkers int, process func(T) (R, error)) *WorkPool[T, R] {
	return &WorkPool[T, R]{
		numWorkers: numWorkers,
		jobs:       make(chan Task[T, R]),
		results:    make(chan Task[T, R]),
		errors:     make(chan Task[T, R]),
		process:    process,
		quit:       make(chan struct{}),
	}
}

func NewAndStartWorkPool[T any, R any](numWorkers int, process func(T) (R, error)) *WorkPool[T, R] {
	workPool := NewWorkPool(numWorkers, process)
	workPool.Start()
	return workPool
}

// Start starts the worker pool
func (wp *WorkPool[T, R]) Start() {
	for w := 1; w <= wp.numWorkers; w++ {
		wp.wg.Add(1)
		go wp.start()
	}
}

func (wp *WorkPool[T, R]) Stop() {
	wp.once.Do(func() {
		close(wp.quit)
		close(wp.jobs)
		wp.wg.Wait()
		close(wp.results)
		close(wp.errors)
	})
}

func (wp *WorkPool[T, R]) AddTask(task Task[T, R]) {
	wp.jobs <- task
}

func (wp *WorkPool[T, R]) Results() <-chan Task[T, R] {
	return wp.results
}

func (wp *WorkPool[T, R]) Errors() <-chan Task[T, R] {
	return wp.errors
}

func (wp *WorkPool[T, R]) start() {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.quit:
			return
		case job := <-wp.jobs:
			result, err := wp.process(job.Input)
			if err != nil {
				job.Err = fmt.Errorf("error processing job id: %s err: %w", job.ID, err)
				wp.errors <- job
			} else {
				job.Result = result
				wp.results <- job
			}
		}
	}
}

func collectResults[T any, R any](wp *WorkPool[T, R]) []Task[T, R] {
	var results []Task[T, R]
	for result := range wp.Results() {
		if result.Err == nil {
			results = append(results, result)
		}
	}
	return results
}

func collectErrors[T any, R any](wp *WorkPool[T, R]) []Task[T, R] {
	var errors []Task[T, R]
	for err := range wp.Errors() {
		if err.Err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (wp *WorkPool[T, R]) CollectResultsAndErrors() ([]Task[T, R], []Task[T, R]) {
	close(wp.jobs)
	wp.wg.Wait()
	wp.Stop()

	results := collectResults(wp)
	errors := collectErrors(wp)
	return results, errors
}
