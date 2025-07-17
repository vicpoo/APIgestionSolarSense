// src/reports/application/put_report_usecase.go
package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type PutReportUseCase struct {
    repo domain.ReportRepository
}

func NewPutReportUseCase(repo domain.ReportRepository) *PutReportUseCase {
    return &PutReportUseCase{repo: repo}
}

func (uc *PutReportUseCase) UpdateReport(ctx context.Context, report *domain.Report) error {
    return uc.repo.Update(ctx, report)
}