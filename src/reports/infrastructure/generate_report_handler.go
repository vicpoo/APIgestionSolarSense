// src/reports/infrastructure/generate_report_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type GenerateReportHandler struct {
    useCase *application.GenerateReportUseCase
}

func NewGenerateReportHandler(useCase *application.GenerateReportUseCase) *GenerateReportHandler {
    return &GenerateReportHandler{useCase: useCase}
}

func (h *GenerateReportHandler) GenerateReport(c *gin.Context) {
    // Obtener email de la URL
    emailFromURL := c.Param("email")
    
    var req domain.GenerateReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validar que el email de la URL coincida con el del JSON
    if emailFromURL != req.RequestedByEmail {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El email en la URL no coincide con el del cuerpo de la solicitud"})
        return
    }

    report, err := h.useCase.GeneratePDFReport(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, report)
}