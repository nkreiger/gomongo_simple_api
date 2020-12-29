package main

import (
	"app/services"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	addr = "8080"
	get  = "GET"
)

func main() {

	// create anew http router
	rtr := mux.NewRouter()

	// health endpoint
	rtr.HandleFunc("/health", services.Health).Methods(get)

	// use go routinue to serve endpoint
	ctx := context.Background()

	GracefullyListenAndServe(ctx, addr, rtr)
}

func GracefullyListenAndServe(ctx context.Context, servePort string, rtr *mux.Router) {
	http.Handle("/", rtr)

	h := &http.Server{
		Addr:              fmt.Sprintf(":%v", servePort),
		Handler:           handlers.CORS()(rtr),
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		log.Printf("serving on port: %v", servePort)
		if err := h.ListenAndServe(); err != nil {
			log.Fatalf("%v", err)
		}
	}()

	// wait for signal to end
	<-sig

	log.Println("Shutting down server...")

	_ = h.Shutdown(ctx)

	log.Println("server gracefully shutdown")
}
