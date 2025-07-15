//api/srcreports/application/report_service.go	

package application

import (
    "context"
    "github.com/vicpoo/apigestion-solar-go/src/reports/domain"

)

type ReportService struct {
    repo domain.ReportRepository
}

func NewReportService(repo domain.ReportRepository) *ReportService {
    return &ReportService{repo: repo}
}

func (s *ReportService) CreateReport(ctx context.Context, report *domain.Report) error {
    return s.repo.Create(ctx, report)
}

func (s *ReportService) GetReport(ctx context.Context, id int) (*domain.Report, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *ReportService) GetUserReports(ctx context.Context, userID int) ([]domain.Report, error) {
    return s.repo.GetByUserID(ctx, userID)
}

func (s *ReportService) DeleteReport(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}