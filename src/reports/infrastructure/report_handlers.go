//api/src/reports/infrastructure/report_handlers.go

package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type ReportHandlers struct {
    service *application.ReportService
}

func NewReportHandlers(service *application.ReportService) *ReportHandlers {
    return &ReportHandlers{service: service}
}

func (h *ReportHandlers) CreateReport(c *gin.Context) {
    var report domain.Report
    if err := c.ShouldBindJSON(&report); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateReport(c.Request.Context(), &report); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, report)
}

func (h *ReportHandlers) GetReport(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    report, err := h.service.GetReport(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
        return
    }
    c.JSON(http.StatusOK, report)
}

func (h *ReportHandlers) GetUserReports(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    reports, err := h.service.GetUserReports(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, reports)
}

func (h *ReportHandlers) DeleteReport(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.service.DeleteReport(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}