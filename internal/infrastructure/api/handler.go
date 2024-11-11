package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"test_task_NT/internal/application"
	"test_task_NT/internal/domain"
)

type Handler struct {
	service *application.BenchmarkService
}

func NewHandler(service *application.BenchmarkService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RunBenchmark(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Received request to run benchmark")
	var req struct {
		DSN         string `json:"dsn"`
		Query       string `json:"query"`
		DurationMS  int    `json:"duration_ms"`
		Concurrency int    `json:"concurrency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Error("Failed to decode request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	logrus.WithFields(logrus.Fields{
		"dsn":         req.DSN,
		"query":       req.Query,
		"duration_ms": req.DurationMS,
		"concurrency": req.Concurrency,
	}).Info("Request payload decoded successfully")

	benchmark := domain.Benchmark{
		Query:       req.Query,
		DurationMS:  req.DurationMS,
		Concurrency: req.Concurrency,
	}

	result, err := h.service.RunBenchmark(r.Context(), req.DSN, benchmark)
	if err != nil {
		logrus.WithError(err).Error("Benchmark execution failed")
		http.Error(w, "Benchmark execution failed", http.StatusInternalServerError)
		return
	}
	logrus.Info("Benchmark executed successfully")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	logrus.Info("Response sent successfully")
}
