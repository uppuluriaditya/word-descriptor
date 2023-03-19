package workerqueue

import (
	"log"
	"sync"
)

type Dispatcher struct {
	// A pool of workers that are registered with the dispatcher
	WorkerPool chan chan Job

	// number of worker pools to be started
	NumWorkers int

	wg *sync.WaitGroup

	workers []Worker
}

func NewDispatcher(numWorkers int) *Dispatcher {
	JobQueue = make(chan Job)
	pool := make(chan chan Job, numWorkers)
	return &Dispatcher{WorkerPool: pool, NumWorkers: numWorkers, wg: &sync.WaitGroup{}}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	d.wg.Add(d.NumWorkers)
	for i := 0; i < d.NumWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.wg, i)
		d.workers = append(d.workers, *worker)
		worker.Start()
	}

	go d.dispatch()
	d.wg.Wait()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				log.Printf("Enqueing the Job: %v\n", job.String())
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Dispatcher) Close() {
	for _, worker := range d.workers {
		worker.Stop()
	}
}
