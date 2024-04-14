package routers

import (
	"github.com/Numostanley/BankingAPI/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func GetTransactionRouters() *chi.Mux {
	transactionRouter := chi.NewRouter()

	transactionRouter.Post("/", handlers.CreateTransactionHandler)
	transactionRouter.Get("/{reference}", handlers.GetTransactionHandler)

	return transactionRouter
}
