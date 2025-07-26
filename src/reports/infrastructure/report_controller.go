// src/reports/infrastructure/report_controller.go
package infrastructure

import "github.com/gin-gonic/gin"

type ReportController struct {
    getHandler    *GetReportHandler
    postHandler   *PostReportHandler
    putHandler    *PutReportHandler
    deleteHandler *DeleteReportHandler
}

func NewReportController(
    getHandler *GetReportHandler,
    postHandler *PostReportHandler,
    putHandler *PutReportHandler,
    deleteHandler *DeleteReportHandler,
) *ReportController {
    return &ReportController{
        getHandler:    getHandler,
        postHandler:   postHandler,
        putHandler:    putHandler,
        deleteHandler: deleteHandler,
    }
}

func (c *ReportController) CreateReport(ctx *gin.Context) {
    c.postHandler.CreateReport(ctx)
}

func (c *ReportController) GetReport(ctx *gin.Context) {
    c.getHandler.GetReport(ctx)
}

func (c *ReportController) GetUserReports(ctx *gin.Context) {
    c.getHandler.GetUserReports(ctx)
}

func (c *ReportController) UpdateReport(ctx *gin.Context) {
    c.putHandler.UpdateReport(ctx)
}

func (c *ReportController) DeleteReport(ctx *gin.Context) {
    c.deleteHandler.DeleteReport(ctx)
}

func (c *ReportController) GetAllReports(ctx *gin.Context) {
    c.getHandler.GetAllReports(ctx)
}

func (c *ReportController) GetReportsByDate(ctx *gin.Context) {
    c.getHandler.GetReportsByDate(ctx)
}