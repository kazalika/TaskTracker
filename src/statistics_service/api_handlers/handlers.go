package statistics_service

import (
	"clickhouse_handlers"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	defaultTopUsersSize = 3
	defaultTopTasksSize = 5
)

func GetTaskStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task_id := mux.Vars(r)["task_id"]

	statistics, err := clickhouse_handlers.GetTaskStatistics(task_id)
	if err != nil {
		err = fmt.Errorf("`GetTaskStatistics failed with message: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statistics["task_id"] = task_id

	encoded, err := json.Marshal(statistics)
	if err != nil {
		err = fmt.Errorf("json result marshal error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encoded)
}

func GetTopTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	parameter := mux.Vars(r)["parameter"]
	top, err := clickhouse_handlers.GetTopTasksByParameter(parameter, defaultTopTasksSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := json.Marshal(top)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encoded)
}

func GetTopUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	top, err := clickhouse_handlers.GetTopUsersByLikes(defaultTopUsersSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := json.Marshal(top)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encoded)
}
