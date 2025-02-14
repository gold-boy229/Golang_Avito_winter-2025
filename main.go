package main

import (
	"MerchShop/controllers"
	"MerchShop/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	// Load Configurations from config.json using Viper
	LoadAppConfig()

	// Initialize Database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Register Routes
	RegisterProductRoutes(router)

	// Start the server
	log.Printf(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

func RegisterProductRoutes(router *mux.Router) {

	router.HandleFunc("/api/auth", controllers.GetAuthTokenHandler).Methods(http.MethodPost)

	router.Handle("/api/buy/{item}", useJWTMiddleware(controllers.BuyItemHandler)).Methods(http.MethodGet)
	router.Handle("/api/sendCoin", useJWTMiddleware(controllers.SendCoinHandler)).Methods(http.MethodPost)
	router.Handle("/api/info", useJWTMiddleware(controllers.GetInfoHandler)).Methods(http.MethodGet)

	// Попробовать сделать что-то через группировку путей:
	// jwt-authentication-golang  file main.go  func initRouter()

	router.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(controllers.MethodNotAllowedHandler)
}

func useJWTMiddleware(func(http.ResponseWriter, *http.Request)) http.Handler {
	return controllers.JWTMiddleware(http.HandlerFunc(controllers.BuyItemHandler))
}
