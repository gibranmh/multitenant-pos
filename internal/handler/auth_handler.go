package handler

import (
	"fmt"
	"multitenant-pos/configs"
	"multitenant-pos/internal/middleware"
	"multitenant-pos/internal/model"
	"multitenant-pos/internal/utils"
	"net/http"
	"time"
)

var Users = map[string]model.User{}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username or password", http.StatusNotAcceptable)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	user := model.User{
		Username:     username,
		Password:     hashedPassword,
		SessionToken: "",
		CSRFToken:    "",
	}

	result := configs.DB.Create(&user)

	if result.Error != nil {
		http.Error(w, "Gagal menyimpan user ke database: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user model.User
	result := configs.DB.First(&user, "username = ?", username)

	if result.Error != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	configs.DB.Save(&user)

	fmt.Fprintln(w, "User logged in successfully")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	user.SessionToken = ""
	user.CSRFToken = ""

	configs.DB.Save(&user)

	fmt.Fprintln(w, "User logged out successfully")
}
