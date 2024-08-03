package main

import (
	"log"

	"github.com/kylejryan/mocument/internal/config"
	"github.com/kylejryan/mocument/internal/db"
	"github.com/kylejryan/mocument/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	docDB, err := db.NewDocDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to DocDB: %v", err)
	}

	handler := handlers.NewHandler(docDB)
	handler.HandleRequests()
}
