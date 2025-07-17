// src/memberships/application/delete_membership_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type DeleteMembershipUseCase struct {
    repo domain.MembershipRepository
}

func NewDeleteMembershipUseCase(repo domain.MembershipRepository) *DeleteMembershipUseCase {
    return &DeleteMembershipUseCase{repo: repo}
}

func (uc *DeleteMembershipUseCase) DeleteMembership(ctx context.Context, userID int) error {
    return uc.repo.Delete(ctx, userID)
}