package handler

import (
	"fmt"
	"multitenant-pos/internal/middleware"
	"net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		sendJSONResponse(w, http.StatusMethodNotAllowed, false, "Invalid Method")
		return
	}

	// Authorize the request using JWT token
	claimsData, err := middleware.Authorize(r)
	if err != nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "Unauthorized: "+err.Error())
		return
	}

	// If authorization is successful, return protected data
	message := fmt.Sprintf("Validation Successful! Your User ID: %d, Your Tenant ID: %s", claimsData.UserID, claimsData.TenantID)
	sendJSONResponse(w, http.StatusOK, true, message)
}
