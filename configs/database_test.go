package configs

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestConnectDB(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Log("Warning: Could not load .env file from test path, trying default system env...")
	}

	ConnectDB()

	if DB == nil {
		t.Error("Failed to connect to database")
	}
}
