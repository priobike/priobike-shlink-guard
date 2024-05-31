package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Result      string `json:"result"`
	IsSuperuser bool   `json:"is_superuser"`
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	// Load username and password from environment variables
	usernames := os.Getenv("USERNAMES")
	passwords := os.Getenv("PASSWORDS")
	// Split by comma
	usernamesList := strings.Split(usernames, ",")
	passwordsList := strings.Split(passwords, ",")

	// Set content type json
	w.Header().Set("Content-Type", "application/json")

	// Read the request body
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Authentication attempt: %+v\n", authRequest)

	// Check if the username and password are correct
	var authResponse AuthResponse
	var authenticated = false

	for i, username := range usernamesList {
		if authRequest.Username == username && authRequest.Password == passwordsList[i] {
			authenticated = true
			break
		}
	}

	if authenticated {
		authResponse = AuthResponse{Result: "allow", IsSuperuser: false}
		w.WriteHeader(http.StatusOK)
	} else {
		authResponse = AuthResponse{Result: "deny", IsSuperuser: false}
		w.WriteHeader(http.StatusUnauthorized)
	}

	json.NewEncoder(w).Encode(authResponse)
}

func main() {
	fmt.Println("Starting server")
	http.HandleFunc("/", authHandler)
	// Get port from env
	port := os.Getenv("PORT")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
