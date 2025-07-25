// api/src/login/domain/entities.go
package domain

import "time"

type AuthResponse struct {
    Success  bool   `json:"success"`
    Message  string `json:"message,omitempty"`
    Error    string `json:"error,omitempty"`
    Token    string `json:"token,omitempty"`
    AuthType string `json:"auth_type,omitempty"`
    UserID   int64  `json:"user_id,omitempty"`
    Email    string `json:"email,omitempty"`
    Username string `json:"username,omitempty"`
    IsAdmin  bool   `json:"is_admin,omitempty"`
}


type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

type User struct {
	ID        int64     `json:"id"`
	UID       string    `json:"uid,omitempty"`
	Email     string    `json:"email"`
	Username  string    `json:"username,omitempty"`
	AuthType  string    `json:"auth_type"`
	IsActive  bool      `json:"is_active"`
	LastLogin time.Time `json:"last_login,omitempty"`
	PhotoURL  string    `json:"photo_url,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}