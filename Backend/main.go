package main

import (
	"backend/db"
	"backend/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db.DbConnection()

	router := mux.NewRouter()
	router.HandleFunc("/stock", routes.GetStockHandler).Methods("GET")
	router.HandleFunc("/stock-recomendation", routes.GetStockRecomendationHandler).Methods("GET")

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CORS_CONFIG_URL")},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsConfig.Handler(router)
	log.Println("Server run: 9950")
	log.Fatal(http.ListenAndServe(":9950", handler))
}
