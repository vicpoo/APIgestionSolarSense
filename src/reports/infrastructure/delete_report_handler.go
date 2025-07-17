// src/reports/infrastructure/delete_report_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
)

type DeleteReportHandler struct {
    useCase *application.DeleteReportUseCase
}

func NewDeleteReportHandler(useCase *application.DeleteReportUseCase) *DeleteReportHandler {
    return &DeleteReportHandler{useCase: useCase}
}

func (h *DeleteReportHandler) DeleteReport(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
        return
    }

    if err := h.useCase.DeleteReport(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}