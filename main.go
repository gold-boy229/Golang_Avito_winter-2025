package main

import (
	"MerchShop/config"
	"MerchShop/controllers"
	"MerchShop/database"
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

	router.HandleFunc("/api/auth", controllers.GetAuthTokenHandler).Methods(http.MethodPost)

	router.Handle("/api/buy/{item}", useJWTMiddleware(controllers.BuyItemHandler)).Methods(http.MethodGet)
	router.Handle("/api/sendCoin", useJWTMiddleware(controllers.SendCoinHandler)).Methods(http.MethodPost)
	router.Handle("/api/info", useJWTMiddleware(controllers.GetInfoHandler)).Methods(http.MethodGet)

	// TODO
	// Попробовать сделать что-то через группировку путей:
	// jwt-authentication-golang  file main.go  func initRouter()

	router.NotFoundHandler = http.HandlerFunc(controllers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(controllers.MethodNotAllowedHandler)
}

func startServer(router *mux.Router) {
	port := config.AppConfig.Port

	log.Print("Starting Server on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

func useJWTMiddleware(func(http.ResponseWriter, *http.Request)) http.Handler {
	return controllers.JWTMiddleware(http.HandlerFunc(controllers.BuyItemHandler))
}
