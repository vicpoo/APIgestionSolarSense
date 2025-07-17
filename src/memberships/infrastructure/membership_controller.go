// src/memberships/infrastructure/membership_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type MembershipController struct {
    getHandler    *GetMembershipHandler
    postHandler   *PostMembershipHandler
    putHandler    *PutMembershipHandler
    deleteHandler *DeleteMembershipHandler
}

func NewMembershipController(
    getHandler *GetMembershipHandler,
    postHandler *PostMembershipHandler,
    putHandler *PutMembershipHandler,
    deleteHandler *DeleteMembershipHandler,
) *MembershipController {
    return &MembershipController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *MembershipController) GetUserMembership(ctx *gin.Context) {
    c.getHandler.GetUserMembership(ctx)
}

func (c *MembershipController) UpgradeToPremium(ctx *gin.Context) {
    c.postHandler.UpgradeToPremium(ctx)
}

func (c *MembershipController) DowngradeToFree(ctx *gin.Context) {
    c.postHandler.DowngradeToFree(ctx)
}

func (c *MembershipController) CreateOrUpdate(ctx *gin.Context) {
    c.putHandler.CreateOrUpdate(ctx)
}

func (c *MembershipController) DeleteMembership(ctx *gin.Context) {
    c.deleteHandler.DeleteMembership(ctx)
}