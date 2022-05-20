package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	serverUrl  = "localhost"
	serverPort = "8000"
)

var r = mux.NewRouter()

func main() {
	if err := run(); err != nil {
		log.Printf("While trying, an error has occurred, %v", err)
		os.Exit(1)
	}
}

func run() error {
	prepareHandlers()

	fmt.Printf("App is starting on port %s...\n", serverPort)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", serverUrl, serverPort), r)
	if err != nil {
		return err
	}

	return nil
}

func prepareHandlers() {
	r = mux.NewRouter()
}
