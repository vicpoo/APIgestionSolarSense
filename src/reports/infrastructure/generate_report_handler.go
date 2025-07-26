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
    var req domain.GenerateReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    report, err := h.useCase.GeneratePDFReport(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, report)
}