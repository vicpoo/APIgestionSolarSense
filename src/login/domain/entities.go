package domain

import "time"

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Token   string `json:"token,omitempty"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"idToken"`
}

type User struct {
	ID        int64
	Email     string
	Username  string
	AuthType  string
	IsActive  bool
	LastLogin time.Time
	PhotoURL  string
	Provider  string
}