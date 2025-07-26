// src/reports/infrastructure/get_report_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
)

type GetReportHandler struct {
    useCase *application.GetReportUseCase
}

func NewGetReportHandler(useCase *application.GetReportUseCase) *GetReportHandler {
    return &GetReportHandler{useCase: useCase}
}

func (h *GetReportHandler) GetReport(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
        return
    }

    report, err := h.useCase.GetReport(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
        return
    }

    c.JSON(http.StatusOK, report)
}

func (h *GetReportHandler) GetUserReports(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    reports, err := h.useCase.GetUserReports(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, reports)
}


func (h *GetReportHandler) GetAllReports(c *gin.Context) {
    reports, err := h.useCase.GetAllReports(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, reports)
}

func (h *GetReportHandler) GetReportsByDate(c *gin.Context) {
    date := c.Param("date")
    reports, err := h.useCase.GetReportsByDate(c.Request.Context(), date)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, reports)
}