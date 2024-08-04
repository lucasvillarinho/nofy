package wpool

import (
	"errors"
	"sync"
)

type Job[T any, R any] struct {
	ID     string
	Input  T
	Result R
	Err    error
}


type Worker[T any, R any] struct {
	job     chan *Job[T, R]
	process func(T) (R, error)
	quit    chan struct{}
	results chan<- *Job[T, R]
	errors  chan<- *Job[T, R]
	wg *sync.WaitGroup
}

type Option[T any, R any] func(*Worker[T, R]) error

func NewWorker[T any, R any](opts ...Option[T, R]) (*Worker[T, R], error) {
	worker := &Worker[T, R]{}
	for _, opt := range opts {
		if err := opt(worker); err != nil {
			return nil, err
		}
	}

	if worker.job == nil {
		return nil, errors.New("job channel is required")
	}
	if worker.process == nil {
		return nil, errors.New("process function is required")
	}
	if worker.quit == nil {
		return nil, errors.New("quit channel is required")
	}
	if worker.wg == nil {
		return nil, errors.New("wait group is required")
	}
	if worker.results == nil {
		return nil, errors.New("results channel is required")
	}
	if worker.errors == nil {
		return nil, errors.New("errors channel is required")
	}

	return worker, nil
}

func WithJobChannel[T any, R any](jobCh chan *Job[T, R]) Option[T, R] {
	return func(w *Worker[T, R]) error {
		if jobCh == nil {
			return errors.New("job channel cannot be nil")
		}
		w.job = jobCh
		return nil
	}
}

func WithProcessFunc[T any, R any](process func(T) (R, error)) Option[T, R] {
	return func(w *Worker[T, R]) error {
		if process == nil {
			return errors.New("process function cannot be nil")
		}
		w.process = process
		return nil
	}
}

func WithQuitChannel[T any, R any](quitCh chan struct{}) Option[T, R] {
	return func(w *Worker[T, R]) error {
		if quitCh == nil {
			return errors.New("quit channel cannot be nil")
		}
		w.quit = quitCh
		return nil
	}
}

func WithWaitGroup[T any, R any](wg *sync.WaitGroup) Option[T, R] {
	return func(w *Worker[T, R]) error {
		if wg == nil {
			return errors.New("wait group cannot be nil")
		}
		w.wg = wg
		return nil
	}
}

func WithResultsChannel[T any, R any](resultsCh chan<- *Job[T, R]) Option[T, R] {
	return func(w *Worker[T, R]) error {
		if resultsCh == nil {
			return errors.New("results channel cannot be nil")
		}
		w.results = resultsCh
		return nil
	}
}

func WithErrorsChannel[T any, R any](errCh chan<- *Job[T, R]) Option[T, R] {
	return func(w *Worker[T, R]) error {
		if errCh == nil {
			return errors.New("errors channel cannot be nil")
		}
		w.errors = errCh
		return nil
	}
}

func (w *Worker[T, R]) start() {
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.quit:
				return
			case job := <-w.job:
				output, err := w.process(job.Input)
				job.Err = err
				job.Result = output
				if err != nil {
					w.errors <- job
				} else {
					w.results <- job
				}
			}
		}
	}()
}
