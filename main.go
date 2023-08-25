package main

import (
	"net/http"
	"os"

	"github.com/harbor-xyz/coding-project/database"
	_ "github.com/harbor-xyz/coding-project/docs"
	"github.com/harbor-xyz/coding-project/server"
)

// @title calendly Backend APIs
// @version 1.0
// @description Calendly Backend APIs
// @BasePath /
func main() {
	database.Init(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	r := server.Init()
	http.ListenAndServe(":8080", r)
}
