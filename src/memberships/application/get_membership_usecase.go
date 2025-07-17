// src/memberships/application/get_membership_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type GetMembershipUseCase struct {
    repo domain.MembershipRepository
}

func NewGetMembershipUseCase(repo domain.MembershipRepository) *GetMembershipUseCase {
    return &GetMembershipUseCase{repo: repo}
}

func (uc *GetMembershipUseCase) GetUserMembership(ctx context.Context, userID int) (*domain.Membership, error) {
    return uc.repo.GetByUserID(ctx, userID)
}