// api/src/reports/domain/report_repository.go
package domain

import (
	"context"
)

type ReportRepository interface {
	Create(ctx context.Context, report *Report) error
	GetByID(ctx context.Context, id int) (*Report, error)
	GetByUserID(ctx context.Context, userID int) ([]Report, error)
	Update(ctx context.Context, report *Report) error
	Delete(ctx context.Context, id int) error
}