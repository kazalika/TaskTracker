/*
 * Пример API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"mongo"
	"net/http"

	"encoding/json"

	"jwt_handlers"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserData(username string, toSaveMap *map[string]string) (int, error) {
	mongoClient := mongo.GetMongoClient()
	collection := mongoClient.Database("users_data").Collection("users")
	var userInformation bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.Background(), filter).Decode(&userInformation)
	if err != nil {
		return http.StatusUnauthorized, errors.New("user not found")
	}
	BSONToStruct(userInformation, toSaveMap)

	fmt.Println("decoded from BSON: ", toSaveMap)
	return 0, nil
}

func GenerateJWT(username string) (string, error) {
	han := jwt_handlers.GetJWTHandlers()
	payload := jwt.MapClaims{
		"username": username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	tokenString, err := token.SignedString(han.JwtPrivate)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}
	return tokenString, nil
}

func structToBSON(data map[string]string) bson.M {
	bsonData := bson.M{}
	for key, value := range data {
		bsonData[key] = value
	}
	return bsonData
}

func BSONToStruct(bsonData bson.M, toSaveMap *map[string]string) {
	data := make(map[string]string)
	for key, value := range bsonData {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		data[strKey] = strValue
	}

	// fmt.Println(data)
	*toSaveMap = data
}

func StoreUserData(username string, data map[string]string) (int, error) {
	mongoClient := mongo.GetMongoClient()
	collection := mongoClient.Database("users_data").Collection("users")

	delete(data, "_id")
	newUserDataBSON := structToBSON(data)

	if CheckIfUserExists(username) {
		filter := bson.D{{Key: "username", Value: username}}
		_, err := collection.ReplaceOne(context.Background(), filter, newUserDataBSON)
		if err != nil {
			fmt.Println("ReplaceOne failed:", err)
			return http.StatusInternalServerError, err
		}
	} else {
		_, err := collection.InsertOne(context.Background(), newUserDataBSON)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = collection.FindOne(context.Background(), newUserDataBSON).Decode(&newUserDataBSON)
		if err != nil {
			fmt.Println("inserted user is not found ", err)
		}
	}
	return 0, nil
}

func CheckIfUserExists(username string) bool {
	mongoClient := mongo.GetMongoClient()
	collection := mongoClient.Database("users_data").Collection("users")
	var userInformation bson.M
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.Background(), filter).Decode(&userInformation)
	return err == nil
}

func HashPassword(password string) string {
	hash := md5.Sum([]byte(password + "SALT"))
	return fmt.Sprintf("%x", hash)
}

func AuthenticatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Декодируем JSON из RequestBody в структуру Credentials
	var creds AuthenticateBody
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedUserData := make(map[string]string)
	code, err := GetUserData(creds.Username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Проверяем, совпадают ли пароли и есть ли в структуре вообще
	storedPassword, ok := storedUserData["password"]
	if !ok || storedPassword != HashPassword(creds.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	tokenString, err := GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storedUserData["token"] = tokenString

	code, err = StoreUserData(creds.Username, storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
	}

	// Устанавливаем токен в Cookie
	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Декодируем JSON из RequestBody в структуру Credentials
	var creds RegisterBody
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if CheckIfUserExists(creds.Username) {
		http.Error(w, "User with this Username does already exist", http.StatusBadRequest)
		return
	}

	// Генерируем токен
	tokenString, err := GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUserData := map[string]string{
		"username": creds.Username,
		"password": HashPassword(creds.Password),
		"token":    tokenString,
	}
	code, err := StoreUserData(creds.Username, newUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
	}

	// Устанавливаем токен в Cookie
	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func MyProfilePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "No cookie", http.StatusUnauthorized)
		return
	}
	// Получаем доступ к ключам JWT
	han := jwt_handlers.GetJWTHandlers()

	// Получаем токен из Cookie
	tokenString := cookie.Value
	payload := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &payload, func(token *jwt.Token) (interface{}, error) {
		return han.JwtPublic, nil
	})

	// Проверяем валидность токена
	if err != nil || !token.Valid {
		http.Error(w, "Invalid jwt token", http.StatusBadRequest)
		return
	}
	// Проверяем, что в токене есть поле username
	if _, ok := payload["username"]; !ok {
		http.Error(w, "Invalid payload in jwt token", http.StatusBadRequest)
		return
	}
	username := payload["username"].(string)

	storedUserData := make(map[string]string)
	code, err := GetUserData(username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
	}

	if v, ok := storedUserData["token"]; !ok || v != tokenString {
		http.Error(w, "The token has expired", http.StatusUnauthorized)
		return
	}

	// Достаем данные, которые нужно обновить
	var creds MyProfileBody
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if creds.FirstName != "" {
		storedUserData["firstName"] = creds.FirstName
	}
	if creds.LastName != "" {
		storedUserData["lastName"] = creds.LastName
	}
	if creds.Birthday != "" {
		storedUserData["birthday"] = creds.Birthday
	}
	if creds.Email != "" {
		storedUserData["email"] = creds.Email
	}
	if creds.PhoneNumber != "" {
		storedUserData["phone"] = creds.PhoneNumber
	}

	code, err = StoreUserData(username, storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
	}

	w.WriteHeader(http.StatusOK)
}
