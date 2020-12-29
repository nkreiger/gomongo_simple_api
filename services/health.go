package services

import (
	"app/mongo"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type status struct {
	Message string `json:"message"`
}

// Health returns the health of the API
var Health = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json")

	log.Println("Health probe received")

	output := status{
		Message: "Health Probe Successful",
	}

	data, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		log.Printf("error encoding json: %s", data)
	}

	_, err = w.Write(data)
	if err != nil {
		log.Printf("error writing response to user: %v", err)
	}
}

var Connection = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx)
	if err != nil {
		writeStatus(&w, http.StatusBadRequest, err.Error())
		return
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("error disconnecting client: %v", err)
		}

		cancel()
	}()

	// write back okay for good connection
	writeStatus(&w, http.StatusOK, nil)
}

func writeStatus(writer *http.ResponseWriter, statusCode int, response interface{}) {
	(*writer).WriteHeader(statusCode)

	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		log.Printf("serialization error: %v", data)
	}

	_, err = (*writer).Write(data)
	if err != nil {
		log.Printf("error writing response %v", err)
	}
}