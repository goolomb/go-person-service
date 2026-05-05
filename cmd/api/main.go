package main

import (
	"fmt"
	"net/http"

	"github.com/goolomb/go-person-service/internal/httpapi"
)

func main() {
	fmt.Println("Starting person service...")

	router := httpapi.NewRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("server stopped: %v\n", err)
	}
}
