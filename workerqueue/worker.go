package workerqueue

import (
	"log"
	"sync"
)

type Worker struct {
	WorkerPool chan chan Job
	JobPool    chan Job
	doneChan   chan bool
	wg         *sync.WaitGroup
	id         int
}

func NewWorker(workerPool chan chan Job, wg *sync.WaitGroup, id int) *Worker {
	return &Worker{
		WorkerPool: workerPool,
		JobPool:    make(chan Job),
		doneChan:   make(chan bool),
		wg:         wg,
		id:         id,
	}
}

func (w Worker) Start() {
	defer w.wg.Done()
	go func() {
		for {
			// register the current worker into the worker pool
			w.WorkerPool <- w.JobPool
			select {
			case job := <-w.JobPool:
				// process the Job
				if err := job.Process(); err != nil {
					log.Printf("wid: %d, Error in processing the Job: %v, err: %v", w.id, job.String(), err)
					// re-enqueue
					log.Printf("wid: %d, Re-enqueing the job: %v", w.id, job.String())
					JobQueue <- job
				}
			case <-w.doneChan:
				// received the stop signal
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.doneChan <- true
	}()
}
