// src/alerts/infrastructure/alert_controller.go
package infrastructure

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/email"
	udomain "github.com/vicpoo/apigestion-solar-go/src/login/domain"
	
	rdomain "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type AlertController struct {
    postHandler   *PostAlertHandler
    getHandler    *GetAlertHandler
    putHandler    *PutAlertHandler
    deleteHandler *DeleteAlertHandler
    emailService  *email.EmailService
    userRepo      udomain.AuthRepository
    reportRepo rdomain.ReportRepository
}

func NewAlertController(
    postHandler *PostAlertHandler,
    getHandler *GetAlertHandler,
    putHandler *PutAlertHandler,
    deleteHandler *DeleteAlertHandler,
    emailService *email.EmailService,
    userRepo udomain.AuthRepository,
    reportRepo rdomain.ReportRepository,
) *AlertController {
    return &AlertController{
        postHandler:   postHandler,
        getHandler:    getHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
        emailService:  emailService,
        userRepo:      userRepo,
        reportRepo:    reportRepo,
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
    
    if request.AdminEmail != "polarsoftsenss@gmail.com" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "El correo del admin debe ser polarsoftsenss@gmail.com"})
        return
    }
    
    exists, err := c.userRepo.EmailExists(ctx.Request.Context(), userEmail)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar el correo del usuario"})
        return
    }
    
    if !exists {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "El correo del usuario no existe en la base de datos"})
        return
    }

    // Obtener información completa del usuario (corregido el error de asignación)
    user, _, err := c.userRepo.FindUserByEmail(ctx.Request.Context(), userEmail)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
        return
    }

    reports, err := c.reportRepo.GetByUserID(ctx.Request.Context(), int(user.ID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener reportes del usuario"})
        return
    }

    // Si no hay reportes, enviar solo el email sin adjunto
    if len(reports) == 0 {
        err = c.emailService.SendAlertEmail(
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
            "message": "Email enviado correctamente (sin adjuntos)",
            "details": gin.H{
                "from":    request.AdminEmail,
                "to":      userEmail,
                "subject": request.Subject,
            },
        })
        return
    }

    // Obtener el último reporte
    latestReport := reports[len(reports)-1]

    // Leer el archivo PDF desde el sistema de archivos
    pdfData, err := os.ReadFile(latestReport.StoragePath)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el archivo PDF"})
        return
    }

    // Crear el adjunto
    attachment := &email.Attachment{
        Data:        pdfData,
        Filename:    latestReport.FileName,
        ContentType: "application/pdf",
    }

    // Enviar el email con el adjunto
    err = c.emailService.SendAlertEmailWithAttachment(
        userEmail,
        request.Subject,
        request.Message,
        attachment,
    )
    
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Email con PDF enviado correctamente",
        "details": gin.H{
            "from":      request.AdminEmail,
            "to":        userEmail,
            "subject":   request.Subject,
            "file_name": latestReport.FileName,
        },
    })
}