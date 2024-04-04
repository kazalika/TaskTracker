package mongo

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongoClient() {
	// Установить строку подключения
	uri := os.Getenv("MONGO_SERVER")
	if uri == "" {
		log.Fatal("You must set your 'MONGO_SERVER' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	// Установить конфигурацию клиента
	clientOptions := options.Client().ApplyURI(uri)

	// Подключиться к серверу MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Ошибка при подключении к MongoDB:", err)
		return
	}
	MongoClient = client

	// Проверить подключение
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Ошибка при проверке подключения к MongoDB:", err)
		return
	}

	fmt.Println("Подключение к MongoDB успешно установлено!")
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}

func CloseMongoClient() {
	MongoClient.Disconnect(context.Background())
}
