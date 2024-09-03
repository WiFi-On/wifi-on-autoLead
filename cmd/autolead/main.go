package main

import (
	"github.com/joho/godotenv"
)

// Айди ростелекома = 52
func main() {
	godotenv.Load("../common/conf/.env")
}