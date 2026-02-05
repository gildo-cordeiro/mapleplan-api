package api

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/middleware"
	"github.com/gorilla/mux"
)

func RegisterRoutes(registry *HandlerRegistry) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", registry.HealthHandler.Health).Methods("GET")

	userRouter := router.PathPrefix("/api/v1/user").Subrouter()
	{
		userRouter.HandleFunc("/profile", registry.UserHandler.UpdateProfile).Methods("PUT")
		userRouter.HandleFunc("/onboarding", registry.UserHandler.UpdateOnboarding).Methods("PUT")
		userRouter.HandleFunc("/search-partners", registry.UserHandler.SearchPartner).Methods("GET")
		userRouter.HandleFunc("/complete-user", registry.UserHandler.GetCompleteUser).Methods("GET")
	}

	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	{
		authRouter.HandleFunc("/signup", registry.AuthHandler.Signup).Methods("POST")
		authRouter.HandleFunc("/login", registry.AuthHandler.Login).Methods("POST")
	}

	// Middlewares here
	router.Use(middleware.AuthenticationMiddleware)

	return router
}
