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
		// Register more specific/static routes first so they don't get captured by the
		// generic {goalID} parameter route (which would treat e.g. "/widget" as a goalID).
		goalRouter.HandleFunc("/widget", registry.GoalHandler.GetWidgetGoals).Methods("GET")
		goalRouter.HandleFunc("/status-counts", registry.GoalHandler.GetStatusCounts).Methods("GET")

		// Routes that operate on a specific goal by its ID. Constrain the {goalID}
		// to a UUID-like pattern to avoid catching static paths.
		goalRouter.HandleFunc("/", registry.GoalHandler.CreateGoal).Methods("POST")
		goalRouter.HandleFunc("/", registry.GoalHandler.GetGoals).Methods("GET")
		goalRouter.HandleFunc("/{goalID:[0-9a-fA-F-]{36}}", registry.GoalHandler.GetGoalByID).Methods("GET")
		goalRouter.HandleFunc("/{goalID:[0-9a-fA-F-]{36}}/status", registry.GoalHandler.UpdateStatus).Methods("PATCH")
		goalRouter.HandleFunc("/{goalID:[0-9a-fA-F-]{36}}", registry.GoalHandler.UpdateGoal).Methods("PUT")
		goalRouter.HandleFunc("/{goalID:[0-9a-fA-F-]{36}}", registry.GoalHandler.DeleteGoal).Methods("DELETE")
	}

	// Middlewares here
	router.Use(middleware.AuthenticationMiddleware)

	return router
}
