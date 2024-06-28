package auth_service

import (
	"context"
	"fmt"
	"kafka_handlers"
	"log"
	"mongo_handlers"
	"net/http"
	"os"

	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	task_servicepb "task_service/proto"
)

var (
	taskServiceGRPCConnection *grpc.ClientConn
	taskServiceClient         task_servicepb.TaskServiceClient
)

func init() {
	var err error
	taskServiceURL, ok := os.LookupEnv("TASK_SERVICE_URL")
	if !ok {
		log.Fatal("No TASK_SERVICE_URL setted but should")
	}
	taskServiceGRPCConnection, err = grpc.NewClient(taskServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	taskServiceClient = task_servicepb.NewTaskServiceClient(taskServiceGRPCConnection)
}

// Authentication handler
//
//	Method: POST
//
//	If password is incorrect returns 401 (Status Unauthorized)
//	If internal error occurred returns 500 (Status Internal Server Error)
//	If request body is not correct returns 400 (Status Bad Request)
func Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Request data decoding
	var creds AuthenticateBody
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedUserData := make(map[string]string)
	code, err := mongo_handlers.GetUserData(creds.Username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	storedPassword, ok := storedUserData["password"]
	if !ok || storedPassword != HashPassword(creds.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateJWTToken(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code, err = mongo_handlers.StoreUserToken(creds.Username, tokenString)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	code, err = mongo_handlers.StoreUserData(creds.Username, storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

// Registration handler
//
//		Method: POST
//
//		If password is incorrect returns 401 (Status Unauthorized)
//		If internal error occurred returns 500 (Status Internal Server Error)
//		If request body is not correct returns 400 (Status Bad Request)
//	 	If user with this username already exists returns 400 (Status Bad Request)
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Request data decoding
	var creds RegisterBody
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mongo_handlers.CheckIfUserExists(creds.Username) {
		http.Error(w, "User with this Username does already exist", http.StatusBadRequest)
		return
	}

	newUserData := map[string]string{
		"username": creds.Username,
		"password": HashPassword(creds.Password),
	}
	code, err := mongo_handlers.StoreUserData(creds.Username, newUserData)
	if err != nil {
		err = fmt.Errorf("error in function `StoreUserData` occurred: %w", err)
		http.Error(w, err.Error(), code)
		return
	}

	tokenString, err := GenerateJWTToken(creds.Username)
	if err != nil {
		err = fmt.Errorf("error in function `GenerateJWTToken` occurred: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	code, err = mongo_handlers.StoreUserToken(creds.Username, tokenString)
	if err != nil {
		err = fmt.Errorf("error in function `StoreUserToken` occurred: %w", err)
		http.Error(w, err.Error(), code)
		return
	}

	// Make Cookie
	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

// UpdateMyProfile handler
//
//	Method: PUT
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If request body is not correct returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	storedUserData := make(map[string]string)
	code, err = mongo_handlers.GetUserData(username, &storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Decoding RequestBody
	var creds ProfileInfo
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Updating by new information
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

	// Store updated info into Mongo
	code, err = mongo_handlers.StoreUserData(username, storedUserData)
	if err != nil {
		http.Error(w, err.Error(), code)
	}

	w.WriteHeader(http.StatusOK)
}

// CreateTask handler
//
//	Method: POST
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If request body is not correct returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Decoding request body
	var creds CreateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send request to Task Service by GRPC
	IDHolder, err := taskServiceClient.CreateTask(context.Background(), &task_servicepb.TaskContent{
		Title:           creds.Title,
		Description:     creds.Description,
		Status:          creds.Status,
		CreatorUsername: username,
	})
	if err != nil {
		err = fmt.Errorf("grpc `CreateTask` request failed with error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send message to Kafka that new task was created so we need to create empty statistics ({likes: 0, views: 0})
	err = kafka_handlers.CreateEmptyStatistics(IDHolder.Id, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write answer into http response
	http_resp := TaskID{TaskID: IDHolder.Id}
	http_resp_bytes, err := json.Marshal(http_resp)
	if err != nil {
		err = fmt.Errorf("json marshaler failed to marshal new task's id but task was already created with id %v. Error message: %w", IDHolder.Id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(http_resp_bytes)
}

// CreateTask handler
//
//	Method: PUT
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If task with this ID doesn't exist or requestor is not an author of the task returns 400 (Status Bad Request)
//	If request body is not correct returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Decoding request body
	var creds UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get variable from URL
	task_id := mux.Vars(r)["task_id"]

	var task task_servicepb.Task
	task.Id = &task_servicepb.TaskID{
		Id: task_id,
	}
	task.Task = &task_servicepb.TaskContent{
		Title:           creds.Title,
		Description:     creds.Description,
		Status:          creds.Status,
		CreatorUsername: username,
	}

	// Send request to Task Service by GRPC
	// If requestor is not author of task then request returns error `NotFound`
	_, err = taskServiceClient.UpdateTask(context.Background(), &task)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf("grpc request `UpdateTask` failed because task with id=%v doesn't exists or requestor is not an author. Error message: %w", task.Id.Id, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = fmt.Errorf("grpc request `UpdateTask` failed with error message: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("Task has been updated succesfully\n"))
}

// DeleteTask handler
//
//	Method: DELETE
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If task with this ID doesn't exist or requestor is not an author of the task returns 400 (Status Bad Request)
//	If request body is not correct returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get variable from URL
	task_id := mux.Vars(r)["task_id"]

	taskID := &task_servicepb.TaskID{Id: task_id}

	// Send request to Task Service by GRPC
	// If requestor is not author of task then request returns error `NotFound`
	_, err = taskServiceClient.DeleteTask(context.Background(), &task_servicepb.RequestByID{
		Id:                taskID,
		RequestorUsername: username,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf("requestor is not an author of the task. GRPC's error message: %w", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = fmt.Errorf("grpc request `DeleteTask` failed with error message: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("Task has been deleted succesfully"))
}

// GetTask handler
//
//	Method: GET
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If request body is not correct returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get variable from URL
	task_id := mux.Vars(r)["task_id"]

	taskID := task_servicepb.TaskID{
		Id: task_id,
	}

	// Send request to Task Service by GRPC
	grpc_resp, err := taskServiceClient.GetTaskById(context.Background(), &task_servicepb.RequestByID{
		Id:                &taskID,
		RequestorUsername: username,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf("task with this id doesn't exist: %w", err)
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			err = fmt.Errorf("grpc request `GetTaskById` failed with error message: %w", err)
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
		err = fmt.Errorf("json marshaler error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(http_resp_bytes)
}

// GetTaskPage handler
//
//	Method: GET
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If request body is not correct returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func GetTaskPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Decoding request body
	var creds TaskListRequest
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send requset to Task Service by GRPC
	grpc_resp, err := taskServiceClient.GetTaskList(context.Background(), &task_servicepb.TaskPageRequest{
		Offset:   creds.Offset,
		PageSize: creds.PageSize,
	})
	if err != nil {
		err = fmt.Errorf("grpc `GetTaskList` failed with message: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal Protobuf to JSON
	marshaler := jsonpb.Marshaler{}
	jsonStr, err := marshaler.MarshalToString(grpc_resp)
	if err != nil {
		err = fmt.Errorf("protobuf to json marshaler failed with message: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(jsonStr))
}

// View handler
//
//	Method: POST
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func View(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get variable from URL
	task_id := mux.Vars(r)["task_id"]

	// Send requset to Task Service by GRPC to get task's author name
	// If task doesn't exists returns error `NotFound`
	grpc_resp, err := taskServiceClient.GetTaskById(context.Background(), &task_servicepb.RequestByID{
		Id:                &task_servicepb.TaskID{Id: task_id},
		RequestorUsername: username, // ?? not neccessary, maybe remove?
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf("task with this id doesn't exist: %w", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = fmt.Errorf("grpc `GetTaskById` failed with message: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send view to Kafka
	err = kafka_handlers.View(username, task_id, grpc_resp.Task.CreatorUsername)
	if err != nil {
		err = fmt.Errorf("`view` message sending caused a error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task has been viewed succesfully\n"))
}

// Like handler
//
//	Method: POST
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func LikeTaskPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get variable from URL
	task_id := mux.Vars(r)["task_id"]

	// Send requset to Task Service by GRPC to get task's author name
	// If task doesn't exists returns error `NotFound`
	grpc_resp, err := taskServiceClient.GetTaskById(context.Background(), &task_servicepb.RequestByID{
		Id:                &task_servicepb.TaskID{Id: task_id},
		RequestorUsername: username, // ?? not neccessary, maybe remove?
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf("task with this id doesn't exist: %w", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = fmt.Errorf("grpc `GetTaskById` failed with message: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send like to Kafka
	err = kafka_handlers.Like(username, task_id, grpc_resp.Task.CreatorUsername)
	if err != nil {
		err = fmt.Errorf("`like` message sending caused a error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task has been liked succesfully\n"))
}

// GetTaskStats handler
//
//	Method: GET
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func GetTaskStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get variable from URL
	task_id := mux.Vars(r)["task_id"]

	// Get statistics for task from Statistics Service
	resp, err := http.Get("http://statistics_service:8090/tasks/" + task_id + "/stats")
	if err != nil {
		err = fmt.Errorf("statistics service cause a error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	CopyResponseToWriter(w, resp)
}

// GetTopTasks handler
//
//	Method: GET
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func GetTopTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if user is authenticated and get his username
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get variable from URL
	parameter := mux.Vars(r)["parameter"]

	// Get top of tasks by parameter from Statistics Service
	resp, err := http.Get("http://statistics_service:8090/top/tasks/" + parameter)
	if err != nil {
		err = fmt.Errorf("statistics service cause a error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	CopyResponseToWriter(w, resp)
}

// GetTopUsers handler
//
//	Method: GET
//
//	If user is not authenticated returns 400 (Status Bad Request)
//	If internal error occurred returns 500 (Status Internal Server Error)
func GetTopUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check if user is authenticated and get his info
	var username string
	code, err := CheckIfUserAuthenticated(r, &username)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Get top of users by likes from Statistics Service
	resp, err := http.Get("http://statistics_service:8090/top/users")
	if err != nil {
		err = fmt.Errorf("statistics service cause a error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	CopyResponseToWriter(w, resp)
}
