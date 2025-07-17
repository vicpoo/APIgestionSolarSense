//src/memberships/domain/membership.go

package domain


import "time"

type Membership struct {
    ID           int       `json:"id"`
    UserID       int       `json:"user_id"`
    Type         string    `json:"type"` // "free", "premium", "admin"
    ExtraStorage bool      `json:"extra_storage"`
    CreatedAt    time.Time `json:"created_at"`
}