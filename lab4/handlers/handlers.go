package handlers

import (
	"encoding/json"
	"lab4/middleware"
	"lab4/models"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to SecureApp!"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.ValidateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username != "admin" || user.Password != "password123" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateJWT(user.Username, "admin")
	if err != nil {
		http.Error(w, "Could not generate JWT", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome, " + username})
}
