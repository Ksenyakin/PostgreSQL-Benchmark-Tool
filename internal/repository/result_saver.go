package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"test_task_NT/internal/domain"
)

func (r *BenchmarkRepository) SaveResult(ctx context.Context, result domain.Result) error {
	filename := "benchmark_results.json"
	log.WithField("filename", filename).Info("Opening file to save benchmark result")

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithError(err).Error("Failed to open file for writing")
		return fmt.Errorf("не удалось открыть файл для записи: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.WithError(err).Warn("Failed to close the file")
		} else {
			log.Info("File closed successfully")
		}
	}()

	data, err := json.Marshal(result)
	if err != nil {
		log.WithError(err).Error("Failed to serialize result to JSON")
		return fmt.Errorf("не удалось сериализовать результат: %w", err)
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		log.WithError(err).Error("Failed to write result to file")
		return fmt.Errorf("не удалось записать результат в файл: %w", err)
	}

	log.WithField("result", result).Info("Benchmark result saved successfully")
	return nil
}
