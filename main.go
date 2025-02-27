package main

import (
	"MerchShop/config"
	"MerchShop/database"
	"MerchShop/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	loadAppConfig()
	initializeDatabase()
	router := initializeRouter()
	registerRoutes(router)

	startServer(router)
}

func loadAppConfig() {
	config.LoadAppConfig()
}

func initializeDatabase() {
	database.Connect(config.AppConfig.ConnectionString)
	database.Migrate()
}

func initializeRouter() *mux.Router {
	return mux.NewRouter().StrictSlash(true)
}

func registerRoutes(router *mux.Router) {

	router.HandleFunc("/api/auth", handlers.GetAuthTokenHandler).Methods(http.MethodPost)

	router.Handle("/api/buy/{item}", useJWTMiddleware(handlers.BuyItemHandler)).Methods(http.MethodGet)
	router.Handle("/api/sendCoin", useJWTMiddleware(handlers.SendCoinHandler)).Methods(http.MethodPost)
	router.Handle("/api/info", useJWTMiddleware(handlers.GetInfoHandler)).Methods(http.MethodGet)

	// TODO
	// Попробовать сделать что-то через группировку путей:
	// jwt-authentication-golang  file main.go  func initRouter()

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)
}

func startServer(router *mux.Router) {
	port := config.AppConfig.Port

	log.Print("Starting Server on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

func useJWTMiddleware(func(http.ResponseWriter, *http.Request)) http.Handler {
	return handlers.JWTMiddleware(http.HandlerFunc(handlers.BuyItemHandler))
}
