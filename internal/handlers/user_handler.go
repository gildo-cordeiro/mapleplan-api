package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gildo-cordeiro/mapleplan-api/internal/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/services"
)

type UserHandler struct {
	UserService services.UserService
}

func (handler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUserDto contract.CreateNewUserDto

	if err := json.NewDecoder(r.Body).Decode(&newUserDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	userId, err := handler.UserService.RegisterUser(newUserDto)

	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"message": errorMessage})
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"id": userId}
	json.NewEncoder(w).Encode(response)
}

func (handler *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginDto contract.LoginDto

	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	user, err := handler.UserService.FindByEmailAndPass(loginDto.Email, loginDto.Password)

	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"message": errorMessage})
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"id":        user.ID,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	}
	json.NewEncoder(w).Encode(response)

}
