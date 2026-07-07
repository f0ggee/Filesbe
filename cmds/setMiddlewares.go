package cmds

import (
	"Kaban/internal/Deliver/Middlewares"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	return mux.NewRouter()
}
func GetPostRouter(newRouter *mux.Router) *mux.Router {
	postRequest := newRouter.PathPrefix("/").Subrouter()
	postRequest.Use(Middlewares.CheckPostRequest)
	return postRequest
}

func GetGetRouter(router *mux.Router) *mux.Router {
	getRequest := router.PathPrefix("/").Subrouter()
	getRequest.Use(Middlewares.CheckerGetRequests)
	return getRequest
}

func SetCheckBotsRouter(router *mux.Router) {
	router.Use(Middlewares.CheckBots)
}

func SetLogging(router *mux.Router) {
	router.Use(Middlewares.Logging)
}
