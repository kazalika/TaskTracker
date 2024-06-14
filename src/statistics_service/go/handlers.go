/*
 * Handlers:
 *		1. getTaskStatistics 	(GET)
 *			<- task identifier 			[string in URL]
 *			-> count of likes and views {likes, views}
 *		2. getTopTasks		 	(GET)
 *			<- parameter for sort 		[string in URL] ("likes" or "views")
 *			-> top 5 tasks by parameter {task_id, likes, views}
 *		3. getTopUsers			(GET)
 *			<- any
 *			-> top 3 users by likes     {author_username, likes}
 */

package statistics_service

import (
	"clickhouse_handlers"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTaskStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task_id := mux.Vars(r)["task_id"]
	likes, err := clickhouse_handlers.GetTaskLikes(task_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	views, err := clickhouse_handlers.GetTaskViews(task_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encoded, err := json.Marshal(map[string]any{
		"task_id": task_id,
		"likes":   likes,
		"views":   views,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encoded)
	w.WriteHeader(http.StatusOK)
}

func GetTop5Tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	parameter := mux.Vars(r)["parameter"]
	top, err := clickhouse_handlers.GetTopTasksByStat(parameter, 5)
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
	w.WriteHeader(http.StatusOK)
}

func GetTop3Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	top, err := clickhouse_handlers.GetTopUsersByLikes(3)
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
	w.WriteHeader(http.StatusOK)
}
