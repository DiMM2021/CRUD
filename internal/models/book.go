package models

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewBook(title string, author string, publishedAt time.Time) *Book {
	return &Book{
		Title:       title,
		Author:      author,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
