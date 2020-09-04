package controllers

import (
	"book-list/models"
	"book-list/repository/book"
	"book-list/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Controllers struct{}

var books []models.Book

func (c *Controllers) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		var books []models.Book
		var cErr models.Error
		br := bookRepo.BookRepo{}
		books, err := br.GetBooks(db, book, books)
		if err != nil {
			cErr.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, cErr)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		utils.SendSuccess(w, &books)
	}
}

// Remove a single book
func (c *Controllers) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var uId uuid.UUID
		params := mux.Vars(r)
		if sId, ok := params["id"]; ok {
			var err error
			uId, err = uuid.Parse(sId)
			if err != nil {
				log.Println(err)
				_, _ = fmt.Fprintf(w, "Bad ID")
				return
			}
			_, err = db.Exec("delete from books where id=$1", uId)
			if err != nil {
				log.Fatal(err)
			}
		}
		_, _ = fmt.Fprintf(w, "Record with ID %q has beed removed from the database.", uId.String())
	}
}

// Update book(s)
func (c *Controllers) UpdateBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newBooks []*models.Book
		_ = json.NewDecoder(r.Body).Decode(&newBooks)
		for _, n := range newBooks {
			_, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4", &n.Title, &n.Author, &n.Year, &n.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
		_ = json.NewEncoder(w).Encode(&newBooks)
	}
}

// Add book(s)
func (c *Controllers) AddBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cErr models.Error
		var newBooks []*models.Book
		_ = json.NewDecoder(r.Body).Decode(&newBooks)
		for _, book := range newBooks {
			if book.Year == 0 || book.Author == "" || book.Title == "" {
				cErr.Message = "Bad Value(s)"
				utils.SendError(w, http.StatusBadRequest, cErr)
				return
			}
			book.ID = uuid.New()
			db.QueryRow("insert into books (id, title, author, year) values ($1, $2, $3, $4)", book.ID, book.Title, book.Author, book.Year)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		utils.SendSuccess(w, &newBooks)
	}
}

// Get a single book
func (c *Controllers) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var matchedBook []models.Book
		var book models.Book
		param := mux.Vars(r)
		if id, ok := param["id"]; ok {
			uId, err := uuid.Parse(id)
			if err != nil {
				_, _ = fmt.Fprintf(w, "Invalid ID")
				return
			}
			rows, err := db.Query("select * from books where id=$1", uId)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
				if err != nil {
					log.Fatal(err)
				}
				matchedBook = append(matchedBook, book)
			}
			_ = json.NewEncoder(w).Encode(&matchedBook)
			return
		}
		_, _ = fmt.Fprintf(w, "No book found by the ID")
	}
}
