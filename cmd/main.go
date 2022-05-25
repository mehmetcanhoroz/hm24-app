package main

import (
	"fmt"
	"github.com/mehmetcanhoroz/hm24-app/handlers"
	"github.com/mehmetcanhoroz/hm24-app/services"
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
	analyseService := services.NewAnalyseService()
	analyseHandler := handlers.NewAnalyseHandler(analyseService)

	r = mux.NewRouter()

	analyseRouter := r.PathPrefix("/analyse").Subrouter()
	analyseRouter.HandleFunc("/html", analyseHandler.GetHtmlContentOfURL)
	analyseRouter.HandleFunc("/html-version", analyseHandler.DetermineHTMLVersion)
	analyseRouter.HandleFunc("/title", analyseHandler.FindHtmlTitleOfURL)
	analyseRouter.HandleFunc("/links", analyseHandler.GetListOfLinkElements)
	analyseRouter.HandleFunc("/hx", analyseHandler.GetCountOfHXElements)
	analyseRouter.HandleFunc("/login-form", analyseHandler.IsThereLoginForm)
}
