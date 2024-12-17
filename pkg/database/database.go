package database

import (
	"crud_golang/config"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewDBConnection() (*sql.DB, error) {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	logrus.WithFields(logrus.Fields{
		"host":    cfg.DBHost,
		"port":    cfg.DBPort,
		"db_name": cfg.DBName,
		"user":    cfg.DBUser,
	}).Info("Attempting to connect to the database")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithError(err).Error("Error opening database connection")
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		logrus.WithError(err).Error("Error connecting to the database")
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	logrus.Info("Successfully connected to the database")

	return db, nil
}
