// src/memberships/domain/membership_repository.go
package domain

import "context"

type MembershipRepository interface {
    GetByUserID(ctx context.Context, userID int) (*Membership, error)
    CreateOrUpdate(ctx context.Context, membership *Membership) error
    UpgradeToPremium(ctx context.Context, userID int) error
    DowngradeToFree(ctx context.Context, userID int) error
    Delete(ctx context.Context, userID int) error
}