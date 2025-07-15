//api/src/sessions/domain/session.go

package domain

import "time"

type Session struct {
    ID           int       `json:"id"`
    UserID       int       `json:"user_id"`
    SessionToken string    `json:"session_token"`
    FirebaseUID  string    `json:"firebase_uid,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    time.Time `json:"expires_at"`
    IsActive     bool      `json:"is_active"`
}