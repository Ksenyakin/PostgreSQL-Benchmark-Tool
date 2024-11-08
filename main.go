package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

func BenchmarkHandler(w http.ResponseWriter, r *http.Request) {
	var cfg Config

	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := benchmarkRPS(cfg)
	if err != nil {
		http.Error(w, "Benchmark execution failed", http.StatusInternalServerError)
		log.Error("Ошибка выполнения бенчмарка: ", err)
		return
	}

	result.Timestamp = time.Now().Format(time.RFC3339)

	if err := saveResult("benchmark_results.json", result); err != nil {
		log.Error("Ошибка при сохранении результата: ", err)
		http.Error(w, "Failed to save result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode result", http.StatusInternalServerError)
		return
	}

	log.Info("Бенчмарк успешно выполнен и результат отправлен")
}

func main() {

	http.HandleFunc("/benchmark", BenchmarkHandler)
	log.Info("Запуск HTTP-сервера на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска HTTP-сервера: ", err)
	}
}
