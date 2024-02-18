package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "paymentApi/cmd/payment/docs"
	"paymentApi/internal/config"
	"paymentApi/internal/handlers"
	"paymentApi/internal/middlewares"
	"paymentApi/pkg/db"
)

// @title API de Pagos
// @description API para el manejo de pagos
// @version 1
// @host localhost:8080
// @BasePath /
func main() {

	config.LoadConfig()

	db.InitDB(config.DSN)

	r := mux.NewRouter()

	// User
	r.HandleFunc("/user", handlers.CreateUser).Methods("POST")

	// Auth
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.Authenticate)
	api.Use(middlewares.AuditLogMiddleware)

	// Rutas protegidas
	api.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	api.HandleFunc("/user/{id}", handlers.UpdateUser).Methods("PUT")
	api.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")

	// Customer
	api.HandleFunc("/customer", handlers.CreateCustomer).Methods("POST")
	api.HandleFunc("/customer/{id}", handlers.GetCustomer).Methods("GET")
	api.HandleFunc("/customer/{id}", handlers.UpdateCustomer).Methods("PUT")
	api.HandleFunc("/customer/{id}", handlers.DeleteCustomer).Methods("DELETE")

	// Merchant
	api.HandleFunc("/merchant", handlers.CreateMerchant).Methods("POST")
	api.HandleFunc("/merchant/{id}", handlers.GetMerchant).Methods("GET")
	api.HandleFunc("/merchant/{id}", handlers.UpdateMerchant).Methods("PUT")
	api.HandleFunc("/merchant/{id}", handlers.DeleteMerchant).Methods("DELETE")

	// Transaction
	api.HandleFunc("/transaction", handlers.CreateTransaction).Methods("POST")
	api.HandleFunc("/transaction/{id}", handlers.GetTransaction).Methods("GET")
	api.HandleFunc("/transaction/{id}", handlers.UpdateTransaction).Methods("PUT")
	api.HandleFunc("/transaction/token/{token}", handlers.GetTransactionByToken).Methods("GET")
	api.HandleFunc("/transaction/process", handlers.ProcessTransaction).Methods("POST")

	log.Println("Servidor escuchando en el puerto 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error al arrancar el servidor: %v", err)
	}
}
