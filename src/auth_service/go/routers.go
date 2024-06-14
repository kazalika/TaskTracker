package auth_service

import (
	"fmt"
	"net/http"
	"strings"

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
		strings.ToUpper("Post"),
		"/authenticate",
		AuthenticatePost,
	},

	Route{
		"RegisterPost",
		strings.ToUpper("Post"),
		"/register",
		RegisterPost,
	},

	Route{
		"MyProfilePut",
		strings.ToUpper("Put"),
		"/myProfile",
		MyProfilePut,
	},

	Route{
		"CreateTaskPost",
		strings.ToUpper("Post"),
		"/tasks/create",
		CreateTaskPost,
	},

	Route{
		"UpdateTaskPut",
		strings.ToUpper("Put"),
		"/tasks/update/{task_id}",
		UpdateTaskPut,
	},

	Route{
		"DeleteTaskPost",
		strings.ToUpper("Delete"),
		"/tasks/delete/{task_id}",
		DeleteTaskDelete,
	},

	Route{
		"GetTaskGet",
		strings.ToUpper("Get"),
		"/tasks/get/{task_id}",
		GetTaskGet,
	},

	Route{
		"GetTaskPageGet",
		strings.ToUpper("Get"),
		"/tasks/getPage",
		GetTaskPageGet,
	},

	Route{
		"ViewTaskPost",
		strings.ToUpper("Post"),
		"/tasks/{task_id}/view",
		ViewTaskPost,
	},

	Route{
		"LikeTaskPost",
		strings.ToUpper("Post"),
		"/tasks/{task_id}/like",
		LikeTaskPost,
	},

	Route{
		"GetTaskStats",
		strings.ToUpper("Get"),
		"/tasks/{task_id}/stats",
		GetTaskStats,
	},

	Route{
		"GetTopTasksGet",
		strings.ToUpper("Get"),
		"/top/tasks/{parameter}",
		GetTopTasksGet,
	},

	Route{
		"GetTopUsers",
		strings.ToUpper("Get"),
		"/top/users",
		GetTopUsersGet,
	},

	Route{
		"AlwaysOKGet",
		strings.ToUpper("Get"),
		"/ok",
		AlwaysOKGet,
	},
}
