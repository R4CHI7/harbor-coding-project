package main

import (
	"os"

	"github.com/harbor-xyz/coding-project/database"
)

func main() {
	database.Init(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
}
