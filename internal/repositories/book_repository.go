package repositories

import (
	"context"
	"crud_golang/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()
var log = logrus.New()

const booksCacheKey = "books_list"

func CreateBook(db *sql.DB, redisClient *redis.Client, book *models.Book) error {
	query := `INSERT INTO books (title, author, published_at, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRow(query, book.Title, book.Author, book.PublishedAt, book.CreatedAt, book.UpdatedAt).Scan(&book.ID)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"query": query,
		}).Error("Error executing query to create book")
		return err
	}

	if err := clearBooksCache(redisClient); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"book_id": book.ID,
		"title":   book.Title,
	}).Info("Book created successfully")
	return nil
}

func GetBookByID(db *sql.DB, id int) (*models.Book, error) {
	var book models.Book
	query := "SELECT id, title, author FROM books WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithFields(logrus.Fields{
				"id": id,
			}).Warn("Book not found")
			return nil, errors.New("book not found")
		}
		log.WithFields(logrus.Fields{
			"error": err,
			"query": query,
		}).Error("Error executing query to get book by ID")
		return nil, err
	}
	return &book, nil
}

func GetBooks(db *sql.DB, redisClient *redis.Client) ([]*models.Book, error) {
	cachedBooks, err := redisClient.Get(ctx, booksCacheKey).Result()
	if err == nil {
		var books []*models.Book
		if err := json.Unmarshal([]byte(cachedBooks), &books); err != nil {
			logrus.WithField("error", err.Error()).Error("Failed to unmarshal cached books")
			return nil, err
		}
		logrus.Info("Returned books from Redis cache")
		return books, nil
	} else if err != redis.Nil {
		logrus.WithField("error", err.Error()).Error("Error fetching books from Redis")
		return nil, err
	}

	rows, err := db.Query(`SELECT id, title, author, published_at, created_at, updated_at FROM books`)
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Error executing query to get books")
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedAt, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			logrus.WithField("error", err.Error()).Error("Error scanning row for book")
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		logrus.WithField("error", err.Error()).Error("Error after iterating over rows")
		return nil, err
	}

	booksJSON, err := json.Marshal(books)
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to marshal books for cache")
		return nil, err
	}

	err = redisClient.Set(ctx, booksCacheKey, booksJSON, 5*time.Minute).Err()
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to set books cache in Redis")
	}

	logrus.Info("Returned books from database and cached in Redis")
	return books, nil
}

func UpdateBook(db *sql.DB, redisClient *redis.Client, book *models.Book, id int) error {
	query := `
        UPDATE books
        SET title = COALESCE($1, title),
            author = COALESCE($2, author),
            published_at = COALESCE($3, published_at),
            updated_at = COALESCE($4, updated_at)
        WHERE id = $5
        RETURNING id
    `
	err := db.QueryRow(query, book.Title, book.Author, book.PublishedAt, book.UpdatedAt, id).Scan(&book.ID)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"query": query,
		}).Error("Error executing query to update book")
		return err
	}

	if err := clearBooksCache(redisClient); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"book_id": book.ID,
		"title":   book.Title,
	}).Info("Book updated successfully")
	return nil
}

func DeleteBook(db *sql.DB, redisClient *redis.Client, id int) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Error executing query to delete book")
		return err
	}

	if err := clearBooksCache(redisClient); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"book_id": id,
	}).Info("Book deleted successfully")
	return nil
}

func clearBooksCache(redisClient *redis.Client) error {
	err := redisClient.Del(ctx, booksCacheKey).Err()
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to delete books cache in Redis")
		return err
	}
	log.Info("Books cache deleted from Redis")
	return nil
}
