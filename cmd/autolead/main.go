package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"wifionAutolead/internal/routes"
)

// Айди ростелекома = 52
func main() {
	godotenv.Load("../common/conf/.env")

	mux := http.NewServeMux()

	routes.RegisterRoutes(mux)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}