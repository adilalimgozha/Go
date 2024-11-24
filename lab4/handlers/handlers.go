package handlers

import (
	"encoding/json"
	"lab4/middleware"
	"lab4/models"
	"net/http"
)

// HomeHandler - Displays a simple message for the home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to SecureApp!"})
}

// LoginHandler - Handles user login and generates JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user input (e.g., username and password)
	if err := models.ValidateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fake user authentication (In real-world, this should check credentials against a DB)
	if user.Username != "admin" || user.Password != "password123" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	token, err := middleware.GenerateJWT(user.Username, "admin")
	if err != nil {
		http.Error(w, "Could not generate JWT", http.StatusInternalServerError)
		return
	}

	// Respond with the JWT token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ProtectedHandler - Handles protected route access (requires JWT)
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Here, we'll use middleware to ensure that only authorized users can access this route.
	username := r.Header.Get("username")

	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Respond with a welcome message and user details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome, " + username})
}
