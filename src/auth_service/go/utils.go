package auth_service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"jwt_handlers"
	"log"
	"mongo"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

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

func GetUserData(username string, mapToLoad *map[string]string) (code int, err error) {
	mongoClient := mongo.GetMongoClient()
	collection := mongoClient.Database("users_data").Collection("users")

	var userInformation bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err = collection.FindOne(context.Background(), filter).Decode(&userInformation)
	if err != nil {
		return http.StatusUnauthorized, errors.New("user not found in registered users database")
	}

	ConvertBSONToStruct(userInformation, mapToLoad)
	return http.StatusOK, nil
}

// Generate JWT token for user
func GenerateJWTToken(username string) (string, error) {
	han := jwt_handlers.GetJWTHandlers()
	payload := jwt.MapClaims{
		"username": username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)

	tokenString, err := token.SignedString(han.JwtPrivate)
	if err != nil {
		err = fmt.Errorf("error while signing token: %w", err)
		return "", err
	}
	return tokenString, nil
}

func CheckIfUserAuthenticated(r *http.Request, username *string, storedUserData *map[string]string) (code int, err error) {
	// Get token from Cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return http.StatusBadRequest, err
	}
	tokenString := cookie.Value

	han := jwt_handlers.GetJWTHandlers()
	payload := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &payload, func(token *jwt.Token) (interface{}, error) {
		return han.JwtPublic, nil
	})

	// Check if token is valid
	if err != nil || !token.Valid {
		return http.StatusBadRequest, errors.New("invalid jwt token")
	}

	// Check if token has neccessary information in payload
	if _, ok := payload["username"]; !ok {
		return http.StatusBadRequest, errors.New("invalid payload in jwt token")
	}
	*username = payload["username"].(string) // ?? is it ok?

	// Get user's information to check if token is relevant
	code, err = GetUserData(*username, storedUserData)
	if err != nil {
		return
	}

	// Check if token is relevant
	if v, ok := (*storedUserData)["token"]; !ok || v != tokenString {
		return http.StatusUnauthorized, errors.New("the token has expired")
	}

	return http.StatusOK, nil
}

func StoreUserData(username string, data map[string]string) (int, error) {
	mongoClient := mongo.GetMongoClient()
	collection := mongoClient.Database("users_data").Collection("users")

	// This meta info will be automatically setted by Mongo, so we should delete it here
	delete(data, "_id")

	newUserDataBSON := ConvertStructToBSON(data)

	if CheckIfUserExists(username) {
		// Replace old info by new one
		filter := bson.D{{Key: "username", Value: username}}
		_, err := collection.ReplaceOne(context.Background(), filter, newUserDataBSON)
		if err != nil {
			err = fmt.Errorf("mongo replace old info by new one is failed with error: %w", err)
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
		// err = collection.FindOne(context.Background(), newUserDataBSON).Decode(&newUserDataBSON)
		// if err != nil {
		// 	fmt.Println("inserted user is not found ", err)
		// }
	}
	return http.StatusOK, nil
}

func CheckIfUserExists(username string) bool {
	mongoClient := mongo.GetMongoClient()
	collection := mongoClient.Database("users_data").Collection("users")
	var userInformation bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.Background(), filter).Decode(&userInformation)
	// If `FindOne` finished incorrectly then user is not found
	// ?? is it ok?
	if err != nil {
		// To see error of it is occured
		log.Println(err.Error())
	}
	return err == nil
}

// Hash given password by md5 + salt
func HashPassword(password string) string {
	hash := md5.Sum([]byte(password + "SALT"))
	return fmt.Sprintf("%x", hash)
}

// Copy response `resp` to response writer `rw`
func CopyResponseToWriter(rw http.ResponseWriter, resp *http.Response) {
	rw.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	rw.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(rw, resp.Body)
	resp.Body.Close()
}
