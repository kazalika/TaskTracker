package main

import (
	"kafka_handlers"
	"log"
	"net/http"

	main_logic "auth_service/main_logic"
	"jwt_handlers"
	"mongo_handlers"
)

func main() {
	log.Printf("[0/3]: Server starting...")

	if err := mongo_handlers.InitMongoClient(); err != nil {
		log.Fatal(err)
	}
	log.Printf("[1/3]: Mongo client initialized")
	defer mongo_handlers.CloseMongoClient()

	if err := jwt_handlers.InitJWTHandlers(); err != nil {
		log.Fatal(err)
	}
	log.Printf("[2/3]: JWT Handlers initialized")

	kafka_handlers.InitKafkaTopics()
	log.Printf("[3/3]: Kafka topics initialized")
	defer kafka_handlers.CloseKafkaTopics()

	router := main_logic.NewRouter()
	log.Println("[Ready] Listen on :8080. You can send requests to main service")
	log.Fatal(http.ListenAndServe(":8080", router))
}
