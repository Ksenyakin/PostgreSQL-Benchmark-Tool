package repository

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type WorkerPool struct {
	numWorkers int
	tasks      chan struct{}
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	log.WithField("num_workers", numWorkers).Info("Creating new worker pool")
	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan struct{}),
	}
}

func (wp *WorkerPool) Start(workerFunc func(id int)) {
	log.WithField("num_workers", wp.numWorkers).Info("Starting worker pool")
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			log.WithField("worker_id", id).Info("Worker started")
			for range wp.tasks {
				log.WithField("worker_id", id).Debug("Worker received task")
				workerFunc(id)
			}
			log.WithField("worker_id", id).Info("Worker stopped")
		}(i)
	}
}

func (wp *WorkerPool) AddTask() {
	log.Debug("Adding task to worker pool")
	wp.tasks <- struct{}{}
}

func (wp *WorkerPool) Stop() {
	log.Info("Stopping worker pool")
	close(wp.tasks)
	wp.wg.Wait()
	log.Info("All workers have stopped")
}
