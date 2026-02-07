package api

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/middleware"
	"github.com/gorilla/mux"
)

func RegisterRoutes(registry *HandlerRegistry) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", registry.HealthHandler.Health).Methods("GET")

	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	{
		authRouter.HandleFunc("/signup", registry.AuthHandler.Signup).Methods("POST")
		authRouter.HandleFunc("/login", registry.AuthHandler.Login).Methods("POST")
	}

	userRouter := router.PathPrefix("/api/v1/user").Subrouter()
	{
		userRouter.HandleFunc("/profile", registry.UserHandler.UpdateProfile).Methods("PUT")
		userRouter.HandleFunc("/onboarding", registry.UserHandler.UpdateOnboarding).Methods("PUT")
		userRouter.HandleFunc("/search-partners", registry.UserHandler.SearchPartner).Methods("GET")
		userRouter.HandleFunc("/complete-user", registry.UserHandler.GetCompleteUser).Methods("GET")
	}

	goalRouter := router.PathPrefix("/api/v1/goals").Subrouter()
	{
		goalRouter.HandleFunc("/", registry.GoalHandler.CreateGoal).Methods("POST")
		goalRouter.HandleFunc("/", registry.GoalHandler.GetGoals).Methods("GET")
		goalRouter.HandleFunc("/widget", registry.GoalHandler.GetWidgetGoals).Methods("GET")
		goalRouter.HandleFunc("/status-counts", registry.GoalHandler.GetStatusCounts).Methods("GET")
		goalRouter.HandleFunc("/{goalID}", registry.GoalHandler.UpdateGoal).Methods("PUT")
		goalRouter.HandleFunc("/{goalID}", registry.GoalHandler.DeleteGoal).Methods("DELETE")
	}

	// Middlewares here
	router.Use(middleware.AuthenticationMiddleware)

	return router
}
