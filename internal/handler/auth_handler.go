package handler

import (
	"encoding/json"
	"multitenant-pos/configs"
	"multitenant-pos/internal/model"
	"multitenant-pos/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Struct Definition for request body in endpoint register and login
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct for API Response to parsing on client side
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Helper for sending JSON response
func sendJSONResponse(w http.ResponseWriter, statusCode int, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Success: success,
		Message: message,
	})
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, false, "Invalid Method")
		return
	}

	// Decode JSON from Frontend to Struct
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "Invalid JSON body format")
		return
	}

	if len(req.Username) < 8 || len(req.Password) < 8 {
		sendJSONResponse(w, http.StatusNotAcceptable, false, "Username or password must be at least 8 characters")
		return
	}

	// Check if username already exists
	var existingUser model.User
	err := configs.DB.First(&existingUser, "username = ?", req.Username).Error
	if err == nil {
		sendJSONResponse(w, http.StatusConflict, false, "Username already exists")
		return
	}

	// Hash the password and save the user to the database
	hashedPassword, _ := utils.HashPassword(req.Password)
	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	result := configs.DB.Create(&user)
	if result.Error != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "Failed to save user: "+result.Error.Error())
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "User registered successfully")
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, false, "Invalid Request Method")
		return
	}

	// Decode JSON from Frontend to Struct
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "Invalid JSON body format")
		return
	}

	// Validate user credentials
	var user model.User
	result := configs.DB.First(&user, "username = ?", req.Username)
	if result.Error != nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "Invalid username or password")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		sendJSONResponse(w, http.StatusUnauthorized, false, "Invalid username or password")
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"tenant_id": user.TenantID,
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign the token with the secret key
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "Failed to generate token")
		return
	}

	// Set the JWT token in an HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
	})

	sendJSONResponse(w, http.StatusOK, true, "User logged in successfully")
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, false, "Invalid Request Method")
		return
	}

	// Clear the JWT token cookie by setting it to an empty value and expiring it immediately
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
	})

	sendJSONResponse(w, http.StatusOK, true, "User logged out successfully")
}
