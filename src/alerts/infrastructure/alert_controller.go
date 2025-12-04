// src/alerts/infrastructure/alert_controller.go
package infrastructure

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/apigestion-solar-go/src/email"
	udomain "github.com/vicpoo/apigestion-solar-go/src/login/domain"

	"github.com/vicpoo/apigestion-solar-go/src/reports/domain"
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

    // Obtener información básica del usuario (sin password hash para usuarios de Google)
    user, _, err := c.userRepo.FindUserByEmail(ctx.Request.Context(), userEmail)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario: " + err.Error()})
        return
    }

    reports, err := c.reportRepo.GetByUserID(ctx.Request.Context(), int(user.ID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener reportes del usuario: " + err.Error()})
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
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar email: " + err.Error()})
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
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el archivo PDF: " + err.Error()})
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
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar email con adjunto: " + err.Error()})
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


func (c *AlertController) CheckSensorAlerts(ctx *gin.Context) {
    userEmail := ctx.Param("userEmail")
    
    // Verificar que el correo del admin sea el correcto
    var request struct {
        AdminEmail string `json:"admin_email"`
    }
    
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if request.AdminEmail != "polarsoftsenss@gmail.com" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "El correo del admin debe ser polarsoftsenss@gmail.com"})
        return
    }
    
    // Verificar que el usuario existe
    exists, err := c.userRepo.EmailExists(ctx.Request.Context(), userEmail)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar el correo del usuario"})
        return
    }
    
    if !exists {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "El correo del usuario no existe en la base de datos"})
        return
    }

    // Mapa de nombres de sensores según su ID
    sensorNames := map[int]string{
        5: "DHT11 (Humedad)",
        6: "BMP280 (Presión)",
        7: "DS18B20 (Temperatura)",
        9: "Otro Sensor", // Actualiza esto según tus necesidades
    }

    // Obtener la fecha actual
    currentDate := time.Now().Format("2006-01-02")
    
    // Obtener los últimos 5 registros de sensor_readings para la fecha actual
    readings, err := c.reportRepo.GetSensorReadingsByDate(ctx.Request.Context(), currentDate)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener lecturas de sensores"})
        return
    }
    
    // Limitar a los últimos 5 registros
    var lastReadings []domain.SensorReading
    if len(readings) > 5 {
        lastReadings = readings[:5]
    } else {
        lastReadings = readings
    }
    
    // Verificar condiciones de alerta
    var alertMessages []string
    
    for _, reading := range lastReadings {
        if reading.Humidity != nil && *reading.Humidity > 80 {
            alertMessages = append(alertMessages, 
                fmt.Sprintf("Alerta de humedad: %.2f%% (mayor que 80%%) - Sensor: %s", 
                    *reading.Humidity, sensorNames[reading.SensorID]))
        }
        
        if reading.Temperature != nil && *reading.Temperature > 35 {
            alertMessages = append(alertMessages, 
                fmt.Sprintf("Alerta de temperatura: %.2f°C (mayor que 35°C) - Sensor: %s", 
                    *reading.Temperature, sensorNames[reading.SensorID]))
        }
        
        if reading.Pressure != nil && *reading.Pressure < 990 {
            alertMessages = append(alertMessages, 
                fmt.Sprintf("Alerta de presión: %.2fhPa (menor que 990hPa) - Sensor: %s", 
                    *reading.Pressure, sensorNames[reading.SensorID]))
        }
    }
    
    // Si no hay alertas, responder sin enviar correo
    if len(alertMessages) == 0 {
        ctx.JSON(http.StatusOK, gin.H{
            "status": "no_alerts",
            "message": "No se detectaron condiciones de alerta",
            "readings": lastReadings,
        })
        return
    }
    
    // Construir el mensaje del correo
    subject := "Alertas de sensores detectadas"
    body := "Se han detectado las siguientes condiciones de alerta:\n\n"
    for _, msg := range alertMessages {
        body += "- " + msg + "\n"
    }
    body += "\nÚltimas lecturas:\n"
    for _, reading := range lastReadings {
        sensorName, exists := sensorNames[reading.SensorID]
        if !exists {
            sensorName = fmt.Sprintf("Sensor %d", reading.SensorID)
        }
        
        body += fmt.Sprintf("- %s: Temp=%.2f°C, Hum=%.2f%%, Pres=%.2fhPa, Registrado: %s\n",
            sensorName,
            safeFloat(reading.Temperature),
            safeFloat(reading.Humidity),
            safeFloat(reading.Pressure),
            reading.RecordedAt.Format("2006-01-02 15:04:05"))
    }
    
    // Enviar el correo
    err = c.emailService.SendAlertEmail(userEmail, subject, body)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Error al enviar el correo de alerta",
            "details": err.Error(),
        })
        return
    }
    
    // Preparar respuesta con nombres de sensores
    var responseReadings []map[string]interface{}
    for _, reading := range lastReadings {
        sensorName, exists := sensorNames[reading.SensorID]
        if !exists {
            sensorName = fmt.Sprintf("Sensor %d", reading.SensorID)
        }
        
        responseReadings = append(responseReadings, map[string]interface{}{
            "sensor_id":   reading.SensorID,
            "sensor_name": sensorName,
            "temperature": reading.Temperature,
            "humidity":    reading.Humidity,
            "pressure":    reading.Pressure,
            "recorded_at": reading.RecordedAt.Format("2006-01-02 15:04:05"),
        })
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "status": "alerts_sent",
        "message": "Alertas enviadas por correo",
        "alerts": alertMessages,
        "readings": responseReadings,
        "user_email": userEmail,
        "admin_email": request.AdminEmail,
    })
}
// Función auxiliar para manejar valores nulos
func safeFloat(f *float64) float64 {
    if f == nil {
        return 0
    }
    return *f
}