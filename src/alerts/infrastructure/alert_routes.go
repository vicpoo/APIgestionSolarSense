// src/alerts/infrastructure/alert_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "github.com/vicpoo/apigestion-solar-go/src/email"
    authinfra "github.com/vicpoo/apigestion-solar-go/src/login/infrastructure"
    reportinfra "github.com/vicpoo/apigestion-solar-go/src/reports/infrastructure"
)

func InitAlertRoutes(router *gin.Engine, emailService *email.EmailService) {
    repo := NewMySQLAlertRepository()
    db := core.GetBD()

    // Crear repositorio de reportes
    reportRepo := reportinfra.NewMySQLReportRepository()

    // Crear casos de uso
    postUseCase := application.NewPostAlertUseCase(repo)
    getUseCase := application.NewGetAlertUseCase(repo)
    putUseCase := application.NewPutAlertUseCase(repo)
    deleteUseCase := application.NewDeleteAlertUseCase(repo)

    // Crear handlers
    postHandler := NewPostAlertHandler(postUseCase)
    getHandler := NewGetAlertHandler(getUseCase)
    putHandler := NewPutAlertHandler(putUseCase)
    deleteHandler := NewDeleteAlertHandler(deleteUseCase)

    // Crear controlador con el userRepo y reportRepo
    controller := NewAlertController(
        postHandler,
        getHandler,
        putHandler,
        deleteHandler,
        emailService,
        authinfra.NewAuthRepository(db),
        reportRepo,
    )

    // Configurar rutas
    alertGroup := router.Group("/api/alerts")
    {
        alertGroup.POST("/", controller.CreateAlert)
        alertGroup.GET("/:id", controller.GetAlert)
        alertGroup.GET("/sensor/:sensor_id", controller.GetSensorAlerts)
        alertGroup.GET("/unsent", controller.GetUnsentAlerts)
        alertGroup.PUT("/:id/mark-sent", controller.MarkAlertAsSent)
        alertGroup.PUT("/:id", controller.UpdateAlert)
        alertGroup.DELETE("/:id", controller.DeleteAlert)
        alertGroup.POST("/test-email/:userEmail", controller.TestEmailAlert)
        alertGroup.POST("/check-alerts/:userEmail", controller.CheckSensorAlerts)
    }
}