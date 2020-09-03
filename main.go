package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
)

// Model
type Book struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Year   int       `json:"year"`
}

var books []Book
var db *sql.DB

func init() {
	err := gotenv.Load()
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANT_SQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgURL)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	r := mux.NewRouter()

	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBooks).Methods("POST")
	r.HandleFunc("/books", updateBooks).Methods("PUT")
	r.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}

// Remove a single book
func removeBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	params := mux.Vars(r)
	if sId, ok := params["id"]; ok {
		uId, err := uuid.Parse(sId)
		if err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(&books)
			return
		}
		for index, book := range books {
			if book.ID == uId {
				books = append(books[:index], books[index+1:]...)
			}
		}
	}
	_ = json.NewEncoder(w).Encode(&books)
}

// Update book(s)
func updateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var newBooks []*Book
	_ = json.NewDecoder(r.Body).Decode(&newBooks)
	for ni, n := range newBooks {
		for oi, o := range books {
			if n.ID == o.ID {
				books[oi] = *newBooks[ni]
				break
			}
		}
	}
	_ = json.NewEncoder(w).Encode(&books)
}

// Add book(s)
func addBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var newBooks []*Book
	_ = json.NewDecoder(r.Body).Decode(&newBooks)
	for _, book := range newBooks {
		book.ID = uuid.New()
		books = append(books, *book)
	}
	_ = json.NewEncoder(w).Encode(&books)
}

// Get a single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	param := mux.Vars(r)
	if id, ok := param["id"]; ok {
		uId, err := uuid.Parse(id)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Invalid ID")
			return
		}
		for _, book := range books {
			if book.ID == uId {
				_ = json.NewEncoder(w).Encode(&book)
				return
			}
		}
	}
	_, _ = fmt.Fprintf(w, "No book found by the ID")
}

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var book Book
	var books []Book

	rows, err := db.Query("select * from books")
	defer rows.Close()
	logFatal(err)

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}
	_ = json.NewEncoder(w).Encode(&books)
}
