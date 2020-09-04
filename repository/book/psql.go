package bookRepo

import (
	"book-list/models"
	"database/sql"
)

type BookRepo struct{}

func (b *BookRepo) GetBooks(db *sql.DB, book models.Book, books []models.Book) ([]models.Book, error) {
	rows, err := db.Query("select * from books")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
