package user

import (
	"encoding/json"
	"net/http"

	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/middleware"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	servicePort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
)

type Handler struct {
	UserService servicePort.UserService
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating user profile goes here
}

func (h *Handler) UpdateOnboarding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized: missing user id", http.StatusUnauthorized)
		return
	}

	var onboardingDto contract.UpdateUserOnboardingDto
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

func (h *Handler) SearchPartner(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetCompleteUser(w http.ResponseWriter, r *http.Request) {
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
