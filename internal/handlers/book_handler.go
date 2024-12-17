package handlers

import (
	"context"
	"crud_golang/internal/models"
	"crud_golang/internal/repositories"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)
}

func checkMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		log.WithFields(logrus.Fields{
			"method":   r.Method,
			"expected": expectedMethod,
		}).Warn("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func handleError(w http.ResponseWriter, message string, statusCode int) {
	log.WithFields(logrus.Fields{
		"error":       message,
		"status_code": statusCode,
	}).Error("Handling error")
	http.Error(w, message, statusCode)
}

func GetBookByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkMethod(w, r, http.MethodGet) {
			return
		}

		// Извлекаем ID из пути, используя mux.Vars
		idStr := r.URL.Path[len("/book/"):]

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			handleError(w, "Invalid id parameter", http.StatusBadRequest)
			log.WithField("id", idStr).Warn("Invalid id parameter")
			return
		}

		book, err := repositories.GetBookByID(db, id)
		if err != nil {
			handleError(w, err.Error(), http.StatusNotFound)
			log.WithFields(logrus.Fields{
				"id":    id,
				"error": err.Error(),
			}).Warn("Book not found")
			return
		}

		if err := json.NewEncoder(w).Encode(book); err != nil {
			handleError(w, "Failed to encode response", http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Failed to encode book")
			return
		}
	}
}

func GetBooksHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if !checkMethod(w, r, http.MethodGet) {
			return
		}

		cacheKey := "books_list"

		cachedBooks, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil && cachedBooks != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			_, err := w.Write([]byte(cachedBooks))
			if err != nil {
				log.WithField("error", err.Error()).Error("Failed to write cached books")
				handleError(w, "Failed to send cached data", http.StatusInternalServerError)
			}
			log.Info("Returned books from Redis cache")
			return
		} else if err != redis.Nil {
			log.WithField("error", err.Error()).Error("Redis error")
			handleError(w, "Failed to get data from cache", http.StatusInternalServerError)
			return
		}

		books, err := repositories.GetBooks(db, redisClient)
		if err != nil {
			handleError(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Error fetching books")
			return
		}

		booksJSON, err := json.Marshal(books)
		if err != nil {
			handleError(w, "Failed to encode response", http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Failed to marshal books")
			return
		}

		err = redisClient.Set(ctx, cacheKey, booksJSON, 5*time.Minute).Err()
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to cache books in Redis")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		_, err = w.Write(booksJSON)
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to send books")
			handleError(w, "Failed to send response", http.StatusInternalServerError)
		}
		log.Info("Returned books from database and cached in Redis")
	}
}

func CreateBookHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkMethod(w, r, http.MethodPost) {
			return
		}

		var newBook models.Book
		if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
			handleError(w, "Invalid data", http.StatusBadRequest)
			log.WithField("error", err.Error()).Warn("Failed to decode request body")
			return
		}

		if err := repositories.CreateBook(db, redisClient, &newBook); err != nil {
			handleError(w, "Failed to create book", http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Failed to create book")
			return
		}

		cacheKey := "books_list"
		err := redisClient.Del(context.Background(), cacheKey).Err()
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to delete cache for books")
		} else {
			log.Info("Cache for books deleted after creation of new book")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newBook); err != nil {
			handleError(w, "Failed to encode response", http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Failed to encode new book")
			return
		}

		log.WithFields(logrus.Fields{
			"book": newBook,
		}).Info("Successfully created new book")
	}
}

func UpdateBookHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			idStr := mux.Vars(r)["id"]
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				handleError(w, "Invalid id parameter", http.StatusBadRequest)
				log.WithField("id", idStr).Warn("Invalid id parameter")
				return
			}

			currentBook, err := repositories.GetBookByID(db, id)
			if err != nil {
				handleError(w, err.Error(), http.StatusNotFound)
				log.WithFields(logrus.Fields{
					"id":    id,
					"error": err.Error(),
				}).Warn("Book not found")
				return
			}

			var updatedBook models.Book
			if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
				handleError(w, "Invalid data", http.StatusBadRequest)
				log.WithField("error", err.Error()).Warn("Failed to decode request body")
				return
			}

			if updatedBook.Title != "" {
				currentBook.Title = updatedBook.Title
			}
			if updatedBook.Author != "" {
				currentBook.Author = updatedBook.Author
			}
			if !updatedBook.PublishedAt.IsZero() {
				currentBook.PublishedAt = updatedBook.PublishedAt
			}

			err = repositories.UpdateBook(db, redisClient, currentBook, id)
			if err != nil {
				handleError(w, err.Error(), http.StatusInternalServerError)
				log.WithField("error", err.Error()).Error("Failed to update book")
				return
			}

			cacheKey := fmt.Sprintf("book:%d", id)
			err = redisClient.Del(context.Background(), cacheKey).Err()
			if err != nil {
				log.WithField("error", err.Error()).Error("Failed to delete book cache from Redis")
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(currentBook); err != nil {
				handleError(w, "Failed to encode response", http.StatusInternalServerError)
				log.WithField("error", err.Error()).Error("Failed to encode updated book")
				return
			}

			log.WithField("id", id).Info("Successfully updated book")
			return
		}

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Warn("Invalid request method for update")
	}
}

func DeleteBookHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkMethod(w, r, http.MethodDelete) {
			return
		}

		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			handleError(w, "Invalid id parameter", http.StatusBadRequest)
			log.WithField("id", idStr).Warn("Invalid id parameter")
			return
		}

		err = repositories.DeleteBook(db, redisClient, id)
		if err != nil {
			handleError(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Failed to delete book")
			return
		}

		cacheKey := "books_list"
		err = redisClient.Del(context.Background(), cacheKey).Err()
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to delete cache for books")
		} else {
			log.Info("Cache for books deleted after deleting a book")
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Book successfully deleted")); err != nil {
			log.WithField("error", err.Error()).Error("Failed to write response")
		}

		log.Info("Successfully deleted book")
	}
}
