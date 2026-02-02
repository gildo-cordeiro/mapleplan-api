package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	servicePort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	jwtutil "github.com/gildo-cordeiro/mapleplan-api/pkg/jwt"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

type Handler struct {
	UserService servicePort.UserService
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUserDto contract.CreateNewUserDto
	if err := json.NewDecoder(r.Body).Decode(&newUserDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	userId, err := h.UserService.RegisterUser(newUserDto)
	if err != nil {
		if errors.Is(err, utils.ErrAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "email já cadastrado"})
			return
		}

		utils.Log.Printf("Signup: RegisterUser error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// build JWT token
	token, err := jwtutil.GenerateToken(userId, 24*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "could not generate token"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id": userId,
		},
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginDto contract.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	user, err := h.UserService.FindByEmailAndPass(loginDto.Email, loginDto.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid credentials"})
		return
	}

	token, err := jwtutil.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "could not generate token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}

type forgotPasswordDto struct {
	Email string `json:"email"`
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto forgotPasswordDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// TODO: chamar o método apropriado do UserService (ex: RequestPasswordReset, SendPasswordReset, etc.)
	// Exemplo hipotético:
	// err := h.UserService.RequestPasswordReset(dto.Email)
	// if err != nil { ... }

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset requested"})
}
