package main

import (
	"log"
	"net/http"

	"github.com/Numostanley/BankingAPI/env"
	"github.com/Numostanley/BankingAPI/internal/db"
	"github.com/Numostanley/BankingAPI/internal/routers"
	"github.com/Numostanley/BankingAPI/internal/utils"
)

func main() {
	enV := env.GetEnv{}
	enV.LoadEnv()

	db.InitDB()
	utils.SeedAccount("extras/accounts.json")

	mainRouter := routers.GetRoutes()

	server := &http.Server{
		Handler: mainRouter,
		Addr:    ":" + enV.PortString,
	}

	log.Printf("Server starting on port %v", enV.PortString)
	log.Fatal(server.ListenAndServe())
}
