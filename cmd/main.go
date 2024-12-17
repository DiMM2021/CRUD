package main

import (
	"crud_golang/internal/routes"
	"crud_golang/pkg/database"
	"net/http"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func main() {
	log := logrus.New()

	db, err := database.NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	log.Info("Successfully connected to the database")

	r := routes.SetupRoutes(db)

	log.Info("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

}
