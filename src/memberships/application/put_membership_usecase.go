// src/memberships/application/put_membership_usecase.go
package application

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/memberships/domain"
)

type PutMembershipUseCase struct {
    repo domain.MembershipRepository
}

func NewPutMembershipUseCase(repo domain.MembershipRepository) *PutMembershipUseCase {
    return &PutMembershipUseCase{repo: repo}
}

func (uc *PutMembershipUseCase) CreateOrUpdate(ctx context.Context, membership *domain.Membership) error {
    // Verificar que el contexto sea de Gin
    if ginCtx, ok := ctx.(*gin.Context); ok {
        // Usar el contexto Gin directamente
        return uc.repo.CreateOrUpdate(ginCtx, membership)
    }
    // Si no es un contexto Gin, usar el contexto original
    return uc.repo.CreateOrUpdate(ctx, membership)
}