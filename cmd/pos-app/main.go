package main

import (
	"fmt"
	"multitenant-pos/internal/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/protected", handler.ProtectedHandler)

	fmt.Println("Server berjalan di port :8080")
	http.ListenAndServe(":8080", nil)
}
