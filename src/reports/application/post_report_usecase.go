// src/reports/application/post_report_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type PostReportUseCase struct {
    repo domain.ReportRepository
}

func NewPostReportUseCase(repo domain.ReportRepository) *PostReportUseCase {
    return &PostReportUseCase{repo: repo}
}

func (uc *PostReportUseCase) CreateReport(ctx context.Context, report *domain.Report) error {
    return uc.repo.Create(ctx, report)
}