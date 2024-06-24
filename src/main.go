package main

import (
	"fmt"
	"net/http"

	"example.SMSService.com/pkg/router"
)

func main() {
	router := router.CreateRouter()
	fmt.Println("Starting the server on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting the server")
	}
}
