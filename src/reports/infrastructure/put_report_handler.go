// src/reports/infrastructure/put_report_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type PutReportHandler struct {
    useCase *application.PutReportUseCase
}

func NewPutReportHandler(useCase *application.PutReportUseCase) *PutReportHandler {
    return &PutReportHandler{useCase: useCase}
}

func (h *PutReportHandler) UpdateReport(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
        return
    }

    var report domain.Report
    if err := c.ShouldBindJSON(&report); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Asegurar que el ID de la URL coincide con el del body
    report.ID = id

    if err := h.useCase.UpdateReport(c.Request.Context(), &report); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}