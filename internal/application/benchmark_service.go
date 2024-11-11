package application

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	"test_task_NT/internal/domain"
	"test_task_NT/internal/repository"
)

type BenchmarkService struct {
	repo domain.BenchmarkRepository
}

func NewBenchmarkService(repo domain.BenchmarkRepository) *BenchmarkService {
	return &BenchmarkService{repo: repo}
}

func (s *BenchmarkService) RunBenchmark(ctx context.Context, dsn string, benchmark domain.Benchmark) (domain.Result, error) {
	logrus.Info("Starting RunBenchmark method")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to the database")
		return domain.Result{}, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.WithError(err).Warn("Failed to close the database connection")
		} else {
			logrus.Info("Database connection closed successfully")
		}
	}()

	repo := repository.NewBenchmarkRepository(db)
	logrus.Info("Repository initialized")

	result, err := repo.RunBenchmark(ctx, benchmark)
	if err != nil {
		logrus.WithError(err).Error("Failed to run benchmark")
		return result, err
	}
	logrus.Info("Benchmark run successfully")

	if err := repo.SaveResult(ctx, result); err != nil {
		logrus.WithError(err).Error("Failed to save benchmark result")
		return result, err
	}
	logrus.Info("Benchmark result saved successfully")

	return result, nil
}
