// src/memberships/application/post_membership_usecase.go
package application

import (
	"context"

	"github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type PostMembershipUseCase struct {
	repo domain.MembershipRepository
}

func NewPostMembershipUseCase(repo domain.MembershipRepository) *PostMembershipUseCase {
	return &PostMembershipUseCase{repo: repo}
}

func (uc *PostMembershipUseCase) UpgradeToPremium(ctx context.Context, userID int) error {
	return uc.repo.UpgradeToPremium(ctx, userID)
}

func (uc *PostMembershipUseCase) DowngradeToFree(ctx context.Context, userID int) error {
	return uc.repo.DowngradeToFree(ctx, userID)
}

// Nuevo método
func (uc *PostMembershipUseCase) FixMissingMemberships(ctx context.Context) error {
	return uc.repo.FixMissingMemberships(ctx)
}