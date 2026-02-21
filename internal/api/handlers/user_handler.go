package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gildo-cordeiro/mapleplan-api/internal/api/middleware"
	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/user/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/ports/services"
)

type UserHandler struct {
	UserService services.UserService
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating user profile goes here
}

func (h *UserHandler) UpdateOnboarding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	var onboardingDto request.UpdateUserOnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&onboardingDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	if err := h.UserService.UpdateOnboarding(r.Context(), userID, onboardingDto); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "onboarding updated successfully"})
}

func (h *UserHandler) SearchPartner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	name := r.URL.Query().Get("query")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "name query parameter is required"})
		return
	}

	results, err := h.UserService.SearchPartnerByName(userID, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func (h *UserHandler) GetCompleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	user, err := h.UserService.GetCompleteUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
