package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", "user=ovch dbname=parserHabrVacancy sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	service := NewMyService(repo)
	handler := NewMyHandler(service)

	r := mux.NewRouter()

	r.HandleFunc("/search", handler.SearchHandler).Methods("GET")
	r.HandleFunc("/get/{id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/list", handler.ListHandler).Methods("GET")
	r.HandleFunc("/delete/{id}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/history", handler.ListSearchHistoryHandler).Methods("GET")
	r.HandleFunc("/history/{id}", handler.DeleteSearchHistoryHandler).Methods("DELETE")

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", r)
}
