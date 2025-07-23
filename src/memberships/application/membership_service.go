// src/memberships/application/membership_service.go
package application

import (
    "context"
    "errors"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
    "golang.org/x/crypto/bcrypt"
)

type MembershipService struct {
    repo domain.MembershipRepository
}

func NewMembershipService(repo domain.MembershipRepository) *MembershipService {
    return &MembershipService{repo: repo}
}

func (s *MembershipService) GetUserMembership(ctx context.Context, userID int) (*domain.Membership, error) {
    return s.repo.GetByUserID(ctx, userID)
}

func (s *MembershipService) UpdateMembership(ctx context.Context, membership *domain.Membership) error {
    if !membership.IsValidType(membership.Type) {
        return errors.New("invalid membership type")
    }
    return s.repo.CreateOrUpdate(ctx, membership)
}

func (s *MembershipService) UpgradeToPremium(ctx context.Context, userID int) error {
    return s.repo.UpgradeToPremium(ctx, userID)
}

func (s *MembershipService) DowngradeToFree(ctx context.Context, userID int) error {
    return s.repo.DowngradeToFree(ctx, userID)
}

func (s *MembershipService) DeleteMembership(ctx context.Context, userID int) error {
    return s.repo.Delete(ctx, userID)
}

func (s *MembershipService) UpdateUser(ctx context.Context, userID int, email, username, password *string) error {
	var passwordHash *string
	
	if password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("could not hash password")
		}
		hashedStr := string(hashed)
		passwordHash = &hashedStr
	}

	return s.repo.UpdateUser(ctx, userID, email, username, passwordHash)
}