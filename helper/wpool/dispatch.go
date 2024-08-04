package wpool

import (
	"sync"
)

type Dispatcher[T any, R any] struct {
	numWorkers int
	job        chan *Job[T, R]
	workers    []*Worker[T, R]
	quit       chan struct{}
	onceStart  sync.Once
	onceStop   sync.Once
	process    func(T) (R, error)
	results    chan *Job[T, R]
	errors     chan *Job[T, R]
	wg         *sync.WaitGroup
}

func NewDispatcher[T any, R any](numWorkers int, process func(T) (R, error)) *Dispatcher[T, R] {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	return &Dispatcher[T, R]{
		numWorkers: numWorkers,
		job:        make(chan *Job[T, R], numWorkers),
		workers:    make([]*Worker[T, R], numWorkers),
		quit:       make(chan struct{}),
		wg:         &sync.WaitGroup{},
		results:    make(chan *Job[T, R], numWorkers),
		errors:     make(chan *Job[T, R], numWorkers),
		process:   process,
	}

}

func (d *Dispatcher[T, R]) Start() {
	d.onceStart.Do(func() {
		for i := 0; i < d.numWorkers; i++ {
			worker := &Worker[T, R]{
				job:     d.job,
				process: d.process,
				quit:    d.quit,
				wg:      d.wg,
			}
			worker.start()
		}
	})
}

func (d *Dispatcher[T, R]) Stop() {
	d.onceStop.Do(func() {
		close(d.quit)
		d.wg.Wait()
		close(d.results) // Close results channel after all workers are done
		close(d.errors)  // Close errors channel after all workers are done
	})
}

func (d *Dispatcher[T, R]) Init() {
	for len(d.workers) < d.numWorkers {
		worker, _:= NewWorker[T, R](
			WithJobChannel[T, R](d.job),
			WithProcessFunc[T, R](d.process),
			WithQuitChannel[T, R](d.quit),
			WithWaitGroup[T, R](d.wg),
			WithResultsChannel[T, R](d.results),
			WithErrorsChannel[T, R](d.errors),
		)
		d.workers = append(d.workers, worker)
	}
}

func (d *Dispatcher[T, R]) AddJob(job *Job[T, R]) {
	select {
	case d.job <- job:
	case <-d.quit:
	}
}

func (d *Dispatcher[T, R]) CollectResults() []*Job[T, R] {
	var results []*Job[T, R]
	for result := range d.results {
		results = append(results, result)
	}
	return results
}

func (d *Dispatcher[T, R]) CollectErrors() []*Job[T, R] {
	var errors []*Job[T, R]
	for err := range d.errors {
		errors = append(errors, err)
	}
	return errors
}