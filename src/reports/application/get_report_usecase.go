// src/reports/application/get_report_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type GetReportUseCase struct {
    repo domain.ReportRepository
}

func NewGetReportUseCase(repo domain.ReportRepository) *GetReportUseCase {
    return &GetReportUseCase{repo: repo}
}

func (uc *GetReportUseCase) GetReport(ctx context.Context, id int) (*domain.Report, error) {
    return uc.repo.GetByID(ctx, id)
}

func (uc *GetReportUseCase) GetUserReports(ctx context.Context, userID int) ([]domain.Report, error) {
    return uc.repo.GetByUserID(ctx, userID)
}
func (uc *GetReportUseCase) GetAllReports(ctx context.Context) ([]domain.Report, error) {
    return uc.repo.GetAllReports(ctx)
}

func (uc *GetReportUseCase) GetReportsByDate(ctx context.Context, date string) ([]domain.Report, error) {
    return uc.repo.GetReportsByDate(ctx, date)
}