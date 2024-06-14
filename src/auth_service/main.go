package main

import (
	"kafka_handlers"
	"log"
	"net/http"

	sw "auth_service/go"
	"jwt_handlers"
	"mongo"
)

func main() {
	log.Printf("Server starting")

	mongo.InitMongoClient()
	log.Printf("Mongo client initialized")

	defer mongo.CloseMongoClient()

	jwt_handlers.InitJWTHandlers()
	log.Printf("JWT Handlers initialized")

	kafka_handlers.InitKafkaConnections()
	log.Printf("Kafka connections initialized")

	defer kafka_handlers.CloseKafkaConnections()

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
