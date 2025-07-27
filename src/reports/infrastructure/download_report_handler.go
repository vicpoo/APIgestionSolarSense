// src/reports/infrastructure/download_report_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
    "github.com/vicpoo/apigestion-solar-go/src/login/domain"
    authinfra "github.com/vicpoo/apigestion-solar-go/src/login/infrastructure"
    "github.com/vicpoo/apigestion-solar-go/src/core"
    "strconv"
)

type DownloadReportHandler struct {
    useCase    *application.GetReportUseCase
    authRepo   domain.AuthRepository
}

func NewDownloadReportHandler(useCase *application.GetReportUseCase) *DownloadReportHandler {
    db := core.GetBD()
    return &DownloadReportHandler{
        useCase:    useCase,
        authRepo:   authinfra.NewAuthRepository(db),
    }
}

func (h *DownloadReportHandler) DownloadReport(c *gin.Context) {
    fileName := c.Param("filename") // Ej: "reporte_2025-07-25.pdf"
    userEmail := c.Param("email")   // Ej: "edwindjll25@gmail.com"

    // 1. Verificar que el usuario existe
    user, err := h.authRepo.GetBasicUserInfo(c.Request.Context(), userEmail)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Usuario no encontrado",
            "details": err.Error(),
        })
        return
    }

    // 2. Buscar el reporte en la BD
    report, err := h.useCase.GetReportByFileName(c.Request.Context(), fileName)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Reporte no encontrado",
            "details": err.Error(),
        })
        return
    }

    // Debug: Mostrar IDs para diagnóstico
    c.Writer.Header().Add("X-Debug-User-ID", strconv.FormatInt(user.ID, 10))
    c.Writer.Header().Add("X-Debug-Report-User-ID", strconv.Itoa(report.UserID))

    // 3. Verificar que el usuario es el dueño del reporte
    if report.UserID != int(user.ID) {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "No tienes permiso para descargar este reporte",
            "details": gin.H{
                "user_id": user.ID,
                "report_user_id": report.UserID,
                "email": userEmail,
            },
        })
        return
    }

    // 4. Descargar el archivo
    c.File(report.StoragePath)
}