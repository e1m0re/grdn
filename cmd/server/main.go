package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/logger"
	"github.com/e1m0re/grdn/internal/srvhandler"
	"github.com/e1m0re/grdn/internal/storage"
)

func main() {
	var serverAddr string

	var verboseMode bool

	flag.StringVar(&serverAddr, "a", "localhost:8080", "address and port to run server")
	flag.BoolVar(&verboseMode, "v", false, "Torn on extended logging mode")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}

	flag.Parse()

	loggerLevel := "info"
	if verboseMode {
		loggerLevel = "debug"
	}

	if err := logger.Initialize(loggerLevel); err != nil {
		fmt.Print(err)
		return
	}

	store := storage.NewMemStorage()
	handler := srvhandler.NewHandler(store)
	router := chi.NewRouter()
	router.Use(logger.RequestLogger)
	router.Route("/", func(r chi.Router) {
		router.Get("/", handler.GetMainPage)
		router.Get("/value/{mType}/{mName}", handler.GetMetricValue)
		router.Post("/update/{mType}/{mName}/{mValue}", handler.UpdateMetric)
	})

	fmt.Println("Running server on ", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
