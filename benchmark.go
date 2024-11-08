package main

import (
	"context"
	"database/sql"
	"runtime"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func benchmarkRPS(cfg Config) (Result, error) {
	log.Info("Connecting to PostgreSQL...")
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return Result{}, err
	}
	defer db.Close()

	log.Info("Connected successfully to PostgreSQL.")
	log.WithFields(log.Fields{
		"query":       cfg.Query,
		"duration":    cfg.DurationMS,
		"concurrency": cfg.Concurrency,
	}).Info("Starting benchmark")

	numWorkers := cfg.Concurrency
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	log.WithField("workers", numWorkers).Info("Number of workers for benchmark")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DurationMS)*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var result Result

	queries := make(chan struct{}, numWorkers)
	responseTimes := make([]float64, 0, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for range queries {
				start := time.Now()
				_, err := db.ExecContext(ctx, cfg.Query)
				duration := float64(time.Since(start).Milliseconds())

				mu.Lock()
				if err != nil {
					log.WithFields(log.Fields{
						"worker_id": id,
						"error":     err,
					}).Error("Error executing query")
					result.Failed++
				} else {
					log.WithField("worker_id", id).Debug("Query executed successfully")
					result.Successful++
					if duration > 0 {
						responseTimes = append(responseTimes, duration)
					}
				}
				result.TotalRequests++
				mu.Unlock()
			}
		}(i)
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				close(queries)
				return
			case <-ticker.C:
				queries <- struct{}{}
			}
		}
	}()

	wg.Wait()

	result.ResponseTimes = responseTimes
	if len(responseTimes) > 0 {
		var total float64
		for _, time := range responseTimes {
			total += time
		}
		result.AverageTimeMS = total / float64(len(responseTimes))
	}

	log.WithFields(log.Fields{
		"total_requests":           result.TotalRequests,
		"successful":               result.Successful,
		"failed":                   result.Failed,
		"average_response_time_ms": result.AverageTimeMS,
	}).Info("Benchmark completed")

	return result, nil
}
