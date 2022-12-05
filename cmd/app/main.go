package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mini-sns-ws/cmd/app/handlers"
	"mini-sns-ws/cmd/app/middlewares"

	"github.com/gorilla/mux"
)

var version string // Application version

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()

	if port == "" {
		port = "8081"
	}

	router.Use(middlewares.LogMiddleware)
	router.HandleFunc("/version", getVersion)
	router.HandleFunc("/time", handlers.GetServerTime).Methods(http.MethodGet)

	log.Printf("version %s listening on port %s", version, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, version)
}
