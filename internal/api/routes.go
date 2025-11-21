package api

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(registry *HandlerRegistry) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", registry.HealthHandler.Health).Methods("GET")

	userRouter := router.PathPrefix("/users").Subrouter()
	{
		userRouter.HandleFunc("", registry.UserHandler.CreateUser).Methods("POST")
		userRouter.HandleFunc("/{id}", registry.UserHandler.LoginUser).Methods("POST")
	}

	// Middlewares here
	// router.Use(loggingMiddleware)

	return router
}
