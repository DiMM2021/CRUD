package routes

import (
	"context"
	"crud_golang/internal/handlers"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("Error connecting to Redis: %v", err)
	}

	err = redisClient.Set(context.Background(), "key", "value", 5*time.Minute).Err()
	if err != nil {
		logrus.Errorf("Error setting key in Redis: %v", err)
	} else {
		logrus.Info("Key set in Redis for testing")
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Добро пожаловать в приложение CRUD</h1>")
	})

	r.HandleFunc("/book/create", handlers.CreateBookHandler(db, redisClient)).Methods("POST")
	r.HandleFunc("/book/update/{id:[0-9]+}", handlers.UpdateBookHandler(db, redisClient)).Methods("PUT")
	r.HandleFunc("/book/delete/{id:[0-9]+}", handlers.DeleteBookHandler(db, redisClient)).Methods("DELETE")
	r.HandleFunc("/book/{id:[0-9]+}", handlers.GetBookByIDHandler(db)).Methods("GET")
	r.HandleFunc("/books", handlers.GetBooksHandler(db, redisClient)).Methods("GET")

	logrus.Info("Routes are successfully set up")

	return r
}
