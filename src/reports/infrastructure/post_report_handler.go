// src/reports/infrastructure/post_report_handler.go
package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type PostReportHandler struct {
    useCase *application.PostReportUseCase
}

func NewPostReportHandler(useCase *application.PostReportUseCase) *PostReportHandler {
    return &PostReportHandler{useCase: useCase}
}

func (h *PostReportHandler) CreateReport(c *gin.Context) {
    var report domain.Report
    if err := c.ShouldBindJSON(&report); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.useCase.CreateReport(c.Request.Context(), &report); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, report)
}