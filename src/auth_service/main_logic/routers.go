package auth_service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"AuthenticatePost",
		"POST",
		"/authenticate",
		Authenticate,
	},

	Route{
		"RegisterPost",
		"POST",
		"/register",
		Register,
	},

	Route{
		"UpdateMyProfile",
		"PUT",
		"/profile",
		UpdateMyProfile,
	},

	Route{
		"CreateTask",
		"POST",
		"/tasks/",
		CreateTask,
	},

	Route{
		"UpdateTask",
		"PUT",
		"/tasks/{task_id}",
		UpdateTask,
	},

	Route{
		"DeleteTaskPost",
		"DELETE",
		"/tasks/{task_id}",
		DeleteTask,
	},

	Route{
		"GetTaskPage",
		"GET",
		"/tasks/page",
		GetTaskPage,
	},

	Route{
		"GetTask",
		"GET",
		"/tasks/{task_id}",
		GetTask,
	},

	Route{
		"View",
		"POST",
		"/tasks/{task_id}/view",
		View,
	},

	Route{
		"LikeTaskPost",
		"POST",
		"/tasks/{task_id}/like",
		LikeTaskPost,
	},

	Route{
		"GetTaskStats",
		"GET",
		"/tasks/{task_id}/stats",
		GetTaskStats,
	},

	Route{
		"GetTopTasksGet",
		"GET",
		"/top/tasks/{parameter}",
		GetTopTasks,
	},

	Route{
		"GetTopUsers",
		"GET",
		"/top/users",
		GetTopUsers,
	},
}
