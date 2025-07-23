// src/memberships/domain/membership.go
package domain

import (
	"errors"
	"time"
)

// Errores personalizados
var (
	ErrInvalidCredentials  = errors.New("email, username and password are required")
	ErrPasswordHashing     = errors.New("could not hash password")
	ErrEmailAlreadyExists  = errors.New("email already in use")
	ErrUserCreationFailed  = errors.New("could not create user")
	ErrMembershipCreationFailed = errors.New("could not create membership")
)

type Membership struct {
    ID           int       `json:"id"`
    UserID       int       `json:"user_id"`
    Type         string    `json:"type"` // "free", "premium", "admin"
    ExtraStorage bool      `json:"extra_storage"`
    CreatedAt    time.Time `json:"created_at"`
}


type MembershipChange struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    OldRole   string    `json:"old_role"`
    NewRole   string    `json:"new_role"`
    ChangedBy int       `json:"changed_by"`
    ChangedAt time.Time `json:"changed_at"`
}

// ValidTypes lista los tipos de membresía permitidos para cambios
func (m *Membership) ValidTypes() []string {
    return []string{"free", "premium"}
}

// IsValidType verifica si un tipo de membresía es válido
func (m *Membership) IsValidType(membershipType string) bool {
    for _, t := range m.ValidTypes() {
        if t == membershipType {
            return true
        }
    }
    return false
}

type UserWithMembership struct {
    ID             int        `json:"id"`
    Email          string     `json:"email"`
    DisplayName    string     `json:"display_name"`
    PhotoURL       *string    `json:"photo_url,omitempty"` // Cambiado a puntero
    Provider       string     `json:"provider"`
    IsActive       bool       `json:"is_active"`
    MembershipType string     `json:"membership_type,omitempty"`
    ExtraStorage   bool       `json:"extra_storage"`
    MembershipSince time.Time `json:"membership_since,omitempty"`
}