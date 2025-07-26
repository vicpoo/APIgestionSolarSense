// src/reports/infrastructure/report_routes.go
package infrastructure

import (
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
)

func InitReportRoutes(router *gin.Engine) {
    repo := NewMySQLReportRepository()

    // Crear casos de uso
    getUseCase := application.NewGetReportUseCase(repo)
    postUseCase := application.NewPostReportUseCase(repo)
    putUseCase := application.NewPutReportUseCase(repo)
    deleteUseCase := application.NewDeleteReportUseCase(repo)
    generateUseCase := application.NewGenerateReportUseCase(repo)

    // Crear handlers
    getHandler := NewGetReportHandler(getUseCase)
    postHandler := NewPostReportHandler(postUseCase)
    putHandler := NewPutReportHandler(putUseCase)
    deleteHandler := NewDeleteReportHandler(deleteUseCase)
    generateHandler := NewGenerateReportHandler(generateUseCase)

    // Crear controlador
    controller := NewReportController(getHandler, postHandler, putHandler, deleteHandler)

    // Configurar rutas
    reportGroup := router.Group("/api/reports")
    {
        reportGroup.POST("/", controller.CreateReport)
        reportGroup.GET("/:id", controller.GetReport)
        reportGroup.GET("/user/:user_id", controller.GetUserReports)
        reportGroup.PUT("/:id", controller.UpdateReport)
        reportGroup.DELETE("/:id", controller.DeleteReport)
        reportGroup.POST("/generate", generateHandler.GenerateReport)
        reportGroup.GET("/", controller.GetAllReports)
        reportGroup.GET("/date/:date", controller.GetReportsByDate)
    }
}