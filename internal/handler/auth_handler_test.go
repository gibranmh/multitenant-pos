package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"multitenant-pos/configs"
	"multitenant-pos/internal/model"
	"multitenant-pos/internal/utils"

	"github.com/joho/godotenv"
)

func setupTestDB(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("Peringatan: Tidak bisa memuat file .env, menggunakan default system env...")
	}

	configs.ConnectDB()

	configs.DB.AutoMigrate(&model.User{})

	configs.DB.Exec("DELETE FROM users")
}

func TestRegisterHandler_Success(t *testing.T) {
	setupTestDB(t)

	data := url.Values{}
	data.Set("username", "budi_sudasono")
	data.Set("password", "secretpassword123")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var user model.User
	result := configs.DB.First(&user, "username = ?", "budi_sudasono")
	if result.Error != nil {
		t.Error("user should be registered and saved inside MySQL database")
	}
}

func TestRegisterHandler_InvalidMethod(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	setupTestDB(t)

	hashedPassword, _ := utils.HashPassword("supersecret99")
	dummyUser := model.User{
		Username: "tonowidodo",
		Password: hashedPassword,
	}
	configs.DB.Create(&dummyUser)

	data := url.Values{}
	data.Set("username", "tonowidodo")
	data.Set("password", "supersecret99")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	response := rr.Result()
	cookies := response.Cookies()

	var hasSessionCookie, hasCSRFCookie bool
	for _, cookie := range cookies {
		if cookie.Name == "session_token" && cookie.HttpOnly {
			hasSessionCookie = true
		}
		if cookie.Name == "csrf_token" && !cookie.HttpOnly {
			hasCSRFCookie = true
		}
	}

	if !hasSessionCookie {
		t.Error("missing or invalid session_token cookie")
	}
	if !hasCSRFCookie {
		t.Error("missing or invalid csrf_token cookie")
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	setupTestDB(t)

	hashedPassword, _ := utils.HashPassword("correct_pass")
	dummyUser := model.User{
		Username: "user_test",
		Password: hashedPassword,
	}
	configs.DB.Create(&dummyUser)

	data := url.Values{}
	data.Set("username", "user_test")
	data.Set("password", "wrong_pass")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}
