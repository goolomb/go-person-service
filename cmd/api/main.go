package main

import (
	"fmt"
	"net/http"

	"github.com/goolomb/go-person-service/internal/httpapi"
	"github.com/goolomb/go-person-service/internal/service"
)

func main() {
	fmt.Println("Starting person service...")

	personService := service.PersonService{}
	router := httpapi.NewRouter(personService)
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("server stopped: %v\n", err)
	}
}
