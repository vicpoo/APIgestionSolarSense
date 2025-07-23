// src/memberships/infrastructure/membership_controller.go
package infrastructure

import (
	

	"github.com/gin-gonic/gin"
)

type MembershipController struct {
	getHandler      *GetMembershipHandler
	postHandler     *PostMembershipHandler
	putHandler      *PutMembershipHandler
	deleteHandler   *DeleteMembershipHandler
	registerHandler *RegisterHandler
}

func NewMembershipController(
	getHandler *GetMembershipHandler,
	postHandler *PostMembershipHandler,
	putHandler *PutMembershipHandler,
	deleteHandler *DeleteMembershipHandler,
	registerHandler *RegisterHandler,
) *MembershipController {
	return &MembershipController{
		getHandler:      getHandler,
		postHandler:     postHandler,
		putHandler:      putHandler,
		deleteHandler:   deleteHandler,
		registerHandler: registerHandler,
	}
}

func (c *MembershipController) RegisterUser(ctx *gin.Context) {
	c.registerHandler.RegisterUser(ctx)
}

func (c *MembershipController) GetUserMembership(ctx *gin.Context) {
	c.getHandler.GetUserMembership(ctx)
}

func (c *MembershipController) GetAllUsers(ctx *gin.Context) {
	c.getHandler.GetAllUsers(ctx)
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

func (c *MembershipController) FixMissingMemberships(ctx *gin.Context) {
	c.postHandler.FixMissingMemberships(ctx) // Ahora este método existe en PostMembershipHandler
}