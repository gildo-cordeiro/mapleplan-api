package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/middleware"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/gorilla/mux"
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
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 3
	}

	results, err := h.GoalService.GetWidgetGoals(r.Context(), userID, intLimit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(results) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (h *GoalHandler) GetStatusCounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	results, err := h.GoalService.GetStatusCounts(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func (h *GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	NewGoalDto := request.CreateGoalRequest{}
	if err := json.NewDecoder(r.Body).Decode(&NewGoalDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := h.GoalService.CreateGoal(r.Context(), NewGoalDto)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Goal created successfully"})
}

func (h *GoalHandler) GetGoals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	results, err := h.GoalService.GetGoals(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(results) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (h *GoalHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	goalID := vars["goalID"]
	var status request.UpdateGoalStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		utils.Log.Errorf("Failed to decode request body: %v", err)
		return
	}

	err := h.GoalService.UpdateStatus(r.Context(), goalID, status.Status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Goal status updated successfully"})
}

func (h *GoalHandler) UpdateGoal(w http.ResponseWriter, r *http.Request) {}

func (h *GoalHandler) DeleteGoal(w http.ResponseWriter, r *http.Request) {}
