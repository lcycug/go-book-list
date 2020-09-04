package main

import (
	"book-list/controllers"
	"book-list/driver"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()
	ctrl := controllers.Controllers{}
	r := mux.NewRouter()

	r.HandleFunc("/books", ctrl.GetBooks(db)).Methods("GET")
	r.HandleFunc("/books/{id}", ctrl.GetBook(db)).Methods("GET")
	r.HandleFunc("/books", ctrl.AddBooks(db)).Methods("POST")
	r.HandleFunc("/books", ctrl.UpdateBooks(db)).Methods("PUT")
	r.HandleFunc("/books/{id}", ctrl.RemoveBook(db)).Methods("DELETE")

	fmt.Println("Server is listening to port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
