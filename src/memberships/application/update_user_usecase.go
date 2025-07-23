package application

import (
	"context"
	"github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type UpdateUserUseCase struct {
	repo domain.MembershipRepository
}

func NewUpdateUserUseCase(repo domain.MembershipRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) UpdateUser(ctx context.Context, userID int, email, username, password *string) error {
	return uc.repo.UpdateUser(ctx, userID, email, username, password)
}