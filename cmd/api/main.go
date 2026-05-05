package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/goolomb/go-person-service/internal/httpapi"
	"github.com/goolomb/go-person-service/internal/service"
	"github.com/goolomb/go-person-service/internal/storage"
)

func main() {
	fmt.Println("Starting person service...")

	if err := os.MkdirAll("./data", 0755); err != nil {
		fmt.Printf("failed to create data directory: %v\n", err)
		return
	}

	db, err := storage.OpenSQLiteDB("./data/app.db")
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		return
	}

	repository := storage.NewPersonRepository(db)
	personService := service.NewPersonService(repository)
	router := httpapi.NewRouter(personService)
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("server stopped: %v\n", err)
	}
}
