package mongo_handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoNotFoundErrorMessage = "mongo: no documents in result"
)

var mongoClient *mongo.Client

func InitMongoClient() error {
	uri := os.Getenv("MONGO_SERVER")
	if uri == "" {
		log.Fatal("You must set your 'MONGO_SERVER' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	// Set client configuration
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to the Mongo
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		err = fmt.Errorf("connection error from mongo: %w", err)
		return err
	}
	mongoClient = client

	// Check connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		err = fmt.Errorf("ping error from mongo: %w", err)
		return err
	}

	fmt.Println("Connection to Mongo is complete!")
	return nil
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}

func CloseMongoClient() {
	mongoClient.Disconnect(context.Background())
}

func GetUserToken(username string) (token string, err error) {
	collection := mongoClient.Database("users_data").Collection("tokens")
	var tokenStruct bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err = collection.FindOne(context.Background(), filter).Decode(&tokenStruct)
	if err != nil {
		err = fmt.Errorf("get user's token from mongo failed with message: %w", err)
		return "", err
	}
	mapToConvert := make(map[string]string)
	ConvertBSONToStruct(tokenStruct, &mapToConvert)

	return mapToConvert["token"], nil
}

func StoreUserToken(username string, token string) (code int, err error) {
	collection := mongoClient.Database("users_data").Collection("tokens")

	tokenBSON := ConvertStructToBSON(map[string]string{
		"username": username,
		"token":    token,
	})

	if CheckIfUserExists(username) {
		// Replace old token by new one
		filter := bson.D{{Key: "username", Value: username}}
		_, err := collection.ReplaceOne(context.Background(), filter, tokenBSON)
		if err != nil {
			err = fmt.Errorf("mongo replace old token by new one failed with error: %w", err)
			return http.StatusInternalServerError, err
		}
	} else {
		// Insert new user's token
		_, err := collection.InsertOne(context.Background(), tokenBSON)
		if err != nil {
			err = fmt.Errorf("mongo insert new user's token failed with error: %w", err)
			return http.StatusInternalServerError, err
		}

		// Check if insert is done correctly
		// !! Comment when not debugging
		// err = collection.FindOne(context.Background(), tokenBSON).Decode(&tokenBSON)
		// if err != nil {
		// 	log.Fatal("function `StoreUserToken`: inserted user is not found ", err)
		// }
	}

	return http.StatusOK, nil
}

func StoreUserData(username string, data map[string]string) (code int, err error) {
	collection := mongoClient.Database("users_data").Collection("users")

	// This meta info will be automatically setted by Mongo, so we should delete it here
	delete(data, "_id")

	newUserDataBSON := ConvertStructToBSON(data)

	if CheckIfUserExists(username) {
		// Replace old info by new one
		filter := bson.D{{Key: "username", Value: username}}
		_, err := collection.ReplaceOne(context.Background(), filter, newUserDataBSON)
		if err != nil {
			err = fmt.Errorf("mongo replace old info by new one failed with error: %w", err)
			return http.StatusInternalServerError, err
		}
	} else {
		// Insert new user's info
		_, err := collection.InsertOne(context.Background(), newUserDataBSON)
		if err != nil {
			err = fmt.Errorf("mongo insert new user's info failed with error: %w", err)
			return http.StatusInternalServerError, err
		}

		// Check if insert is done correctly
		// !! Comment when not debugging
		// err = collection.FindOne(context.Background(), newUserDataBSON).Decode(&newUserDataBSON)
		// if err != nil {
		// 	log.Fatal("function `StoreUserData`: inserted user is not found ", err)
		// }
	}
	return http.StatusOK, nil
}

func GetUserData(username string, mapToLoad *map[string]string) (code int, err error) {
	collection := mongoClient.Database("users_data").Collection("users")

	var userInformation bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err = collection.FindOne(context.Background(), filter).Decode(&userInformation)
	// If `FindOne` finished incorrectly then user is not found
	if err != nil {
		if err.Error() != mongoNotFoundErrorMessage {
			log.Println("function `GetUserData`, method `FindOne` returned unexpected error: ", err.Error(), ". service will think that user is not found")
		}
		return http.StatusUnauthorized, errors.New("user not found in registered users database")
	}

	ConvertBSONToStruct(userInformation, mapToLoad)
	return http.StatusOK, nil
}

func CheckIfUserExists(username string) bool {
	collection := mongoClient.Database("users_data").Collection("tokens")
	var userToken bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.Background(), filter).Decode(&userToken)
	// If `FindOne` finished incorrectly then user is not found
	if err != nil {
		if err.Error() != mongoNotFoundErrorMessage {
			log.Println("function `CheckIfUserExists`, method `FindOne` returned unexpected error: ", err.Error(), ". service will think that user is not found")
		}
	}
	return err == nil
}

func ConvertStructToBSON(data map[string]string) bson.M {
	bsonData := bson.M{}
	for key, value := range data {
		bsonData[key] = value
	}
	return bsonData
}

func ConvertBSONToStruct(bsonData bson.M, toSaveMap *map[string]string) {
	for key, value := range bsonData {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		(*toSaveMap)[strKey] = strValue
	}
}
