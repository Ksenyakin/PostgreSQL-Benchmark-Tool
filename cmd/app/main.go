package main

import (
	"log"
	"net/http"
	"test_task_NT/internal/application"
	"test_task_NT/internal/infrastructure/api"
	"test_task_NT/internal/infrastructure/db"
	"test_task_NT/internal/repository"
	"test_task_NT/internal/utils"
)

func main() {
	utils.SetupLogger()

	dbConn, err := db.NewPostgresConnection("your_dsn_here")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConn.Close()

	repo := repository.NewBenchmarkRepository(dbConn)
	service := application.NewBenchmarkService(repo)
	handler := api.NewHandler(service)

	http.HandleFunc("/benchmark", handler.RunBenchmark)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
