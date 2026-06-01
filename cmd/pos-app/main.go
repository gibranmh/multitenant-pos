package main

import (
	"fmt"
	"multitenant-pos/configs"
	"multitenant-pos/internal/handler"
	"multitenant-pos/internal/model"
	"net/http"
)

func main() {
	configs.ConnectDB()
	fmt.Println("Database terhubung dengan sukses!")

	configs.DB.AutoMigrate(&model.User{})
	fmt.Println("Migrasi database sukses!")

	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/protected", handler.ProtectedHandler)

	fmt.Println("Server berjalan di port :8080")
	http.ListenAndServe(":8080", nil)
}
