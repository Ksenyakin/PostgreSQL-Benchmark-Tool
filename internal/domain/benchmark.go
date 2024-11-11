package domain

import "context"

type Benchmark struct {
	Query       string
	DurationMS  int
	Concurrency int
}

type Result struct {
	TotalRequests   int       `json:"total_requests"`
	Successful      int       `json:"successful"`
	Failed          int       `json:"failed"`
	AverageTimeMS   float64   `json:"average_time_ms"`
	MinResponseTime float64   `json:"min_response_time_ms"`
	MaxResponseTime float64   `json:"max_response_time_ms"`
	ResponseTimes   []float64 `json:"response_times,omitempty"`
	Timestamp       string    `json:"timestamp"`
}

type BenchmarkRepository interface {
	RunBenchmark(ctx context.Context, benchmark Benchmark) (Result, error)
	SaveResult(ctx context.Context, result Result) error
}
