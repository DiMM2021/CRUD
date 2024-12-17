CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    published_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);