//src/memberships/application/membership_service.go

package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
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