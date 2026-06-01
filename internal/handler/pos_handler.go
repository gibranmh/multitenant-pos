package handler

import (
	"fmt"
	"multitenant-pos/configs"
	"multitenant-pos/internal/middleware"
	"multitenant-pos/internal/model"
	"net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")

	var user model.User
	result := configs.DB.First(&user, "username = ?", username)

	if result.Error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := middleware.Authorize(r, user.SessionToken, user.CSRFToken); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Validation Successful! Welcome to the protected area, %s!", username)
}
