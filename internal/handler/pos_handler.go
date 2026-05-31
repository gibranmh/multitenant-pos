package handler

import (
	"fmt"
	"multitenant-pos/internal/middleware"
	"net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")

	user, ok := Users[username]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := middleware.Authorize(r, user.SessionToken, user.CSRFToken); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Validation Successful! Welcome to the protected area, %s!", username)
}
