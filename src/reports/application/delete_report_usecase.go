// src/reports/application/delete_report_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type DeleteReportUseCase struct {
    repo domain.ReportRepository
}

func NewDeleteReportUseCase(repo domain.ReportRepository) *DeleteReportUseCase {
    return &DeleteReportUseCase{repo: repo}
}

func (uc *DeleteReportUseCase) DeleteReport(ctx context.Context, id int) error {
    return uc.repo.Delete(ctx, id)
}