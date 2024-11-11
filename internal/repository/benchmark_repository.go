package repository

import (
	"context"
	"database/sql"
	"math"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"test_task_NT/internal/domain"
)

type BenchmarkRepository struct {
	db *sql.DB
}

func NewBenchmarkRepository(db *sql.DB) *BenchmarkRepository {
	return &BenchmarkRepository{db: db}
}

func (r *BenchmarkRepository) RunBenchmark(ctx context.Context, benchmark domain.Benchmark) (domain.Result, error) {
	numWorkers := benchmark.Concurrency
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	log.WithField("workers", numWorkers).Info("Starting benchmark")

	var result domain.Result
	var mu sync.Mutex

	totalDuration := 0.0
	minDuration, maxDuration := math.MaxFloat64, 0.0

	wp := NewWorkerPool(numWorkers)

	workerFunc := func(id int) {
		start := time.Now()
		_, err := r.db.ExecContext(ctx, benchmark.Query)
		duration := float64(time.Since(start).Nanoseconds()) / 1e6

		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			log.WithFields(log.Fields{
				"worker_id": id,
				"error":     err,
			}).Error("Error executing query")
			result.Failed++
		} else {
			result.Successful++
			result.TotalRequests++
			totalDuration += duration
			if duration > 0 && duration < minDuration {
				minDuration = duration
			}
			if duration > maxDuration {
				maxDuration = duration
			}
		}
	}

	log.Info("Initializing worker pool")
	wp.Start(workerFunc)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(benchmark.DurationMS)*time.Millisecond)
	defer cancel()

	go func() {
		for {
			select {
			case <-timeoutCtx.Done():
				log.Info("Benchmark duration reached, stopping tasks")
				return
			default:
				select {
				case wp.tasks <- struct{}{}:
					log.Debug("Task sent to worker pool")
				case <-timeoutCtx.Done():
					return
				}
			}
		}
	}()

	log.Info("Waiting for benchmark completion")
	<-timeoutCtx.Done()
	wp.Stop()

	if result.Successful > 0 {
		result.AverageTimeMS = totalDuration / float64(result.Successful)
	}
	if minDuration == math.MaxFloat64 {
		minDuration = 0
	}
	result.MinResponseTime = minDuration
	result.MaxResponseTime = maxDuration
	result.Timestamp = time.Now().Format(time.RFC3339)

	log.WithFields(log.Fields{
		"total_requests":           result.TotalRequests,
		"successful":               result.Successful,
		"failed":                   result.Failed,
		"average_response_time_ms": result.AverageTimeMS,
		"min_response_time_ms":     result.MinResponseTime,
		"max_response_time_ms":     result.MaxResponseTime,
	}).Info("Benchmark completed")

	return result, nil
}
