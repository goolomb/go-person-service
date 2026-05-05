package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		fmt.Printf("failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	db, err := storage.OpenSQLiteDB(cfg.DBPath)
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		os.Exit(1)
	}

	repository := storage.NewPersonRepository(db)
	personService := service.NewPersonService(repository)
	router := httpapi.NewRouter(personService)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server stopped: %v\n", err)
		os.Exit(1)
	}
}
