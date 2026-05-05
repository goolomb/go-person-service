package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goolomb/go-person-service/internal/config"
	"github.com/goolomb/go-person-service/internal/httpapi"
	"github.com/goolomb/go-person-service/internal/service"
	"github.com/goolomb/go-person-service/internal/storage"
)

func main() {
	fmt.Println("Starting person service...")

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		fmt.Printf("failed to create data directory: %v\n", err)
		return
	}

	db, err := storage.OpenSQLiteDB(cfg.DBPath)
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		return
	}

	repository := storage.NewPersonRepository(db)
	personService := service.NewPersonService(repository)
	router := httpapi.NewRouter(personService)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		fmt.Printf("server stopped: %v\n", err)
	}
}
