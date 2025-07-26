package infrastructure

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/reports/application"
)

type DownloadReportHandler struct {
    useCase *application.GetReportUseCase
}

func NewDownloadReportHandler(useCase *application.GetReportUseCase) *DownloadReportHandler {
    return &DownloadReportHandler{useCase: useCase}
}

func (h *DownloadReportHandler) DownloadReport(c *gin.Context) {
    fileName := c.Param("filename") // Ej: "reporte_2025-07-25.pdf"

    // Buscar el reporte en la BD
    report, err := h.useCase.GetReportByFileName(c.Request.Context(), fileName)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
        return
    }

    // Descargar el archivo
    c.File(report.StoragePath)
}