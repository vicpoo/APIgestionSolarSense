// src/alerts/infrastructure/alert_controller.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/email"
    "net/http"
)

type AlertController struct {
    postHandler   *PostAlertHandler
    getHandler    *GetAlertHandler
    putHandler    *PutAlertHandler
    deleteHandler *DeleteAlertHandler
    emailService  *email.EmailService
}

func NewAlertController(
    postHandler *PostAlertHandler,
    getHandler *GetAlertHandler,
    putHandler *PutAlertHandler,
    deleteHandler *DeleteAlertHandler,
    emailService *email.EmailService,
) *AlertController {
    return &AlertController{
        postHandler:   postHandler,
        getHandler:    getHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
        emailService:  emailService,
    }
}

func (c *AlertController) CreateAlert(ctx *gin.Context) {
    c.postHandler.CreateAlert(ctx)
}

func (c *AlertController) GetAlert(ctx *gin.Context) {
    c.getHandler.GetAlert(ctx)
}

func (c *AlertController) GetSensorAlerts(ctx *gin.Context) {
    c.getHandler.GetSensorAlerts(ctx)
}

func (c *AlertController) GetUnsentAlerts(ctx *gin.Context) {
    c.getHandler.GetUnsentAlerts(ctx)
}

func (c *AlertController) MarkAlertAsSent(ctx *gin.Context) {
    c.postHandler.MarkAlertAsSent(ctx)
}

func (c *AlertController) UpdateAlert(ctx *gin.Context) {
    c.putHandler.UpdateAlert(ctx)
}

func (c *AlertController) DeleteAlert(ctx *gin.Context) {
    c.deleteHandler.DeleteAlert(ctx)
}

// Nuevo método para probar emails
func (c *AlertController) TestEmailAlert(ctx *gin.Context) {
    userEmail := ctx.Param("userEmail")
    
    var request struct {
        AdminEmail string `json:"admin_email"`
        Subject    string `json:"subject"`
        Message    string `json:"message"`
        AlertType  string `json:"alert_type"`
    }
    
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Validar que el admin_email sea polarsoftsenss@gmail.com
    if request.AdminEmail != "polarsoftsenss@gmail.com" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "El correo del admin debe ser polarsoftsens@gmail.com"})
        return
    }
    
    // Enviar el email
    err := c.emailService.SendAlertEmail(
        userEmail,
        request.Subject,
        request.Message,
    )
    
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Email enviado correctamente",
        "details": gin.H{
            "from":    request.AdminEmail,
            "to":      userEmail,
            "subject": request.Subject,
        },
    })
}