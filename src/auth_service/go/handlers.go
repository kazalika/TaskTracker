package auth_service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"kafka_handlers"
	"log"
	"mongo"
	"net/http"

	"encoding/json"

	"jwt_handlers"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	task_servicepb "task_service/proto"
)

var (
	grpc_conn   *grpc.ClientConn
	grpc_client task_servicepb.TaskServiceClient
)

func init() {
	var err error
	grpc_conn, err = grpc.NewClient("dns:///task_service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	grpc_client = task_servicepb.NewTaskServiceClient(grpc_conn)
}

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

func CheckIfUserAuthenticated(r *http.Request, username *string, storedUserData *map[string]string) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return http.StatusUnauthorized, err
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
		return http.StatusBadRequest, errors.New("invalid jwt token")
	}
	// Проверяем, что в токене есть поле username
	if _, ok := payload["username"]; !ok {
		return http.StatusBadRequest, errors.New("invalid payload in jwt token")
	}
	*username = payload["username"].(string)

	code, err := GetUserData(*username, storedUserData)
	if err != nil {
		return code, err
	}

	if v, ok := (*storedUserData)["token"]; !ok || v != tokenString {
		return http.StatusUnauthorized, errors.New("the token has expired")
	}

	return 200, nil
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

	var username string
	storedUserData := make(map[string]string)

	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
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

func CreateTaskPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	var creds CreateTaskBody
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpc_resp, err := grpc_client.CreateTask(context.Background(), &task_servicepb.TaskContent{
		Title:           creds.Title,
		Description:     creds.Description,
		Status:          creds.Status,
		CreatorUsername: username,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Some internal error (with postgres, for example)")
		return
	}

	err = kafka_handlers.CreateEmptyStatistics(grpc_resp.Id, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http_resp := TaskID{TaskID: grpc_resp.Id}

	http_resp_bytes, err := json.Marshal(http_resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(http_resp_bytes)
	w.WriteHeader(http.StatusOK)
}

func UpdateTaskPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	var creds UpdateTaskBody
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task_id := mux.Vars(r)["task_id"]

	var task task_servicepb.Task
	task.Id = &task_servicepb.TaskID{Id: task_id}
	task.Task = &task_servicepb.TaskContent{
		Title:           creds.Title,
		Description:     creds.Description,
		Status:          creds.Status,
		CreatorUsername: username,
	}

	_, err = grpc_client.UpdateTask(context.Background(), &task)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("Task has been updated succesfully"))
	w.WriteHeader(http.StatusOK)
}

func DeleteTaskDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	task_id := mux.Vars(r)["task_id"]
	taskID := &task_servicepb.TaskID{Id: task_id}

	_, err = grpc_client.DeleteTask(context.Background(), &task_servicepb.RequestByID{
		Id:                taskID,
		RequestorUsername: username,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("Task has been deleted succesfully"))
	w.WriteHeader(http.StatusOK)
}

func GetTaskGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	task_id := mux.Vars(r)["task_id"]
	taskID := task_servicepb.TaskID{
		Id: task_id,
	}

	grpc_resp, err := grpc_client.GetTaskById(context.Background(), &task_servicepb.RequestByID{
		Id:                &taskID,
		RequestorUsername: username,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http_resp := TaskContent{
		Title:       grpc_resp.Task.Title,
		Description: grpc_resp.Task.Description,
		Status:      grpc_resp.Task.Status,
	}

	http_resp_bytes, err := json.Marshal(http_resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(http_resp_bytes)
	w.WriteHeader(http.StatusOK)
}

func GetTaskPageGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	var creds TaskListRequest
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpc_resp, err := grpc_client.GetTaskList(context.Background(), &task_servicepb.TaskPageRequest{
		Offset:   creds.Offset,
		PageSize: creds.PageSize,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	marshaler := jsonpb.Marshaler{}
	jsonStr, err := marshaler.MarshalToString(grpc_resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(jsonStr))
	w.WriteHeader(http.StatusOK)
}

/*
 * Send a message about a view into the broker (Kafka)
 */
func ViewTaskPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	task_id := mux.Vars(r)["task_id"]

	grpc_resp, err := grpc_client.GetTaskById(context.Background(), &task_servicepb.RequestByID{
		Id:                &task_servicepb.TaskID{Id: task_id},
		RequestorUsername: username,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = kafka_handlers.View(username, task_id, grpc_resp.Task.CreatorUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
 * Send a message about a like into the broker (Kafka)
 */
func LikeTaskPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	task_id := mux.Vars(r)["task_id"]

	grpc_resp, err := grpc_client.GetTaskById(context.Background(), &task_servicepb.RequestByID{
		Id:                &task_servicepb.TaskID{Id: task_id},
		RequestorUsername: username,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = kafka_handlers.Like(username, task_id, grpc_resp.Task.CreatorUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
 *  Always returns 200 OK
 */
func AlwaysOKGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func pipeReq(rw http.ResponseWriter, resp *http.Response) {
	rw.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	rw.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(rw, resp.Body)
	resp.Body.Close()
}

func GetTaskStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	task_id := mux.Vars(r)["task_id"]

	resp, err := http.Get("http://statistics_service:8090/tasks/" + task_id + "/stats")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pipeReq(w, resp)
}

func GetTopTasksGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	parameter := mux.Vars(r)["parameter"]

	resp, err := http.Get("http://statistics_service:8090/top/tasks/" + parameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pipeReq(w, resp)
}

func GetTopUsersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	storedUserData := make(map[string]string)
	code, err := CheckIfUserAuthenticated(r, &username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	resp, err := http.Get("http://statistics_service:8090/top/users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
	w.WriteHeader(resp.StatusCode)
}
