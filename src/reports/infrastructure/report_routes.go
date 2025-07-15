//api/src/reports/infrastructure/report_routes.go

package infrastructure

import (
    "github.com/gin-gonic/gin"
 
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
  
)

func InitReportRoutes(router *gin.Engine) {
    repo := NewMySQLReportRepository()
    service := application.NewReportService(repo)
    handlers := NewReportHandlers(service)

    reportGroup := router.Group("/api/reports")
    {
        reportGroup.POST("/", handlers.CreateReport)
        reportGroup.GET("/:id", handlers.GetReport)
        reportGroup.GET("/user/:user_id", handlers.GetUserReports)
        reportGroup.DELETE("/:id", handlers.DeleteReport)
    }
}