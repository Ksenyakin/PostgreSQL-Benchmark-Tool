package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func NewPostgresConnection(dsn string) (*sql.DB, error) {
	logrus.Info("Attempting to open a new PostgreSQL connection")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logrus.WithError(err).Error("Failed to open PostgreSQL connection")
		return nil, err
	}
	logrus.Info("PostgreSQL connection opened successfully")
	return db, nil
}
