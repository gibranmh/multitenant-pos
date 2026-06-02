package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsData struct {
	UserID   uint
	TenantID string
}

func Authorize(r *http.Request) (*ClaimsData, error) {
	// Get JWT token from cookie
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		return nil, errors.New("Cookie Not Found")
	}

	tokenString := cookie.Value

	// Parse JWT Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate signing method only HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("Token Invalid")
	}

	// Extract claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Failed to extract")
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if float64(time.Now().Unix()) > exp {
			return nil, errors.New("Token Expired")
		}
	}

	// Extract user ID and tenant ID from claims
	var userID uint
	if sub, ok := claims["sub"].(float64); ok {
		userID = uint(sub)
	}

	var tenantID string
	if tID, ok := claims["tenant_id"].(string); ok {
		tenantID = tID
	}

	// Return the extracted claims data
	return &ClaimsData{
		UserID:   userID,
		TenantID: tenantID,
	}, nil
}
