package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/e1m0re/grdn/internal/srvhandler"
	"github.com/e1m0re/grdn/internal/storage"
)

func main() {
	var serverAddr string

	flag.StringVar(&serverAddr, "a", "localhost:8080", "address and port to run server")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}

	flag.Parse()

	store := storage.NewMemStorage()
	handler := srvhandler.NewHandler(store)
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/", handler.GetMainPage)
		router.Get("/value/{mType}/{mName}", handler.GetMetricValue)
		router.Post("/update/{mType}/{mName}/{mValue}", handler.UpdateMetric)
	})

	fmt.Println("Running server on ", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
