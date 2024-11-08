package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Result struct {
	TotalRequests int       `json:"total_requests"`
	Successful    int       `json:"successful"`
	Failed        int       `json:"failed"`
	AverageTimeMS float64   `json:"average_time_ms"`
	ResponseTimes []float64 `json:"response_times"`
	Timestamp     string    `json:"timestamp"`
}

func saveResult(filename string, result Result) error {
	var results []Result

	if _, err := os.Stat(filename); err == nil {

		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("не удалось прочитать файл с результатами: %w", err)
		}

		// Декодируем существующие данные в слайс результатов
		if err := json.Unmarshal(data, &results); err != nil {
			return fmt.Errorf("не удалось декодировать существующие результаты: %w", err)
		}
	}

	results = append(results, result)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось кодировать результаты: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("не удалось записать результаты в файл: %w", err)
	}

	return nil
}
