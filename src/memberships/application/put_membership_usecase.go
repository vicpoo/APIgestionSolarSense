// src/memberships/application/put_membership_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type PutMembershipUseCase struct {
    repo domain.MembershipRepository
}

func NewPutMembershipUseCase(repo domain.MembershipRepository) *PutMembershipUseCase {
    return &PutMembershipUseCase{repo: repo}
}

func (uc *PutMembershipUseCase) CreateOrUpdate(ctx context.Context, membership *domain.Membership) error {
    return uc.repo.CreateOrUpdate(ctx, membership)
}