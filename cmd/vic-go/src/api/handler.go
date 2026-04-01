package api

import (
	"encoding/json"
	"net/http"
)

// APIHandler handles all HTTP requests
// This is a critical handler - needs tests
type APIHandler struct {
	auth AuthService
}

// AuthService interface for dependency injection
type AuthService interface {
	Authenticate(username, password string) (*User, error)
	ValidateToken(token string) bool
}

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginRequest handles login requests
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse handles login responses
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// Login endpoint for user authentication
// Critical API endpoint - must have tests
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual login logic
	// This should:
	// 1. Parse request body
	// 2. Call auth service
	// 3. Generate JWT
	// 4. Return response

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// This is a placeholder - real implementation would call auth service
	token := "placeholder-token"
	user := &User{
		ID:    "1",
		Name:  req.Username,
		Email: req.Username + "@example.com",
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthCheck endpoint for service health
// Another critical endpoint
func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewAPIHandler(auth AuthService) *APIHandler {
	return &APIHandler{auth: auth}
}