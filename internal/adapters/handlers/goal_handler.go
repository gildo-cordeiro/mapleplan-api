package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/middleware"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
)

type GoalHandler struct {
	GoalService services.GoalService
}

func (h *GoalHandler) GetWidgetGoals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "3" // default limit
	}

	results, err := h.GoalService.GetWidgetGoals(userID, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func (h *GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {}

func (h *GoalHandler) GetGoals(w http.ResponseWriter, r *http.Request) {}

func (h *GoalHandler) UpdateGoal(w http.ResponseWriter, r *http.Request) {}

func (h *GoalHandler) DeleteGoal(w http.ResponseWriter, r *http.Request) {}
