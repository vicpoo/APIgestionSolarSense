// src/memberships/domain/membership.go
package domain

import "time"

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