// src/workers/alert_worker.go
package workers

import (
	"context"
	"log"
	"time"
	"github.com/vicpoo/apigestion-solar-go/src/alerts/application"
)

type AlertWorker struct {
	alertService *application.AlertService
	interval     time.Duration
}

func NewAlertWorker(alertService *application.AlertService, interval time.Duration) *AlertWorker {
	return &AlertWorker{
		alertService: alertService,
		interval:     interval,
	}
}

func (w *AlertWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processUnsentAlerts(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *AlertWorker) processUnsentAlerts(ctx context.Context) {
	alerts, err := w.alertService.GetUnsent(ctx) // Cambiado de GetUnsentAlerts a GetUnsent
	if err != nil {
		log.Printf("Error getting unsent alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		// Reutilizamos CreateAlert que ya tiene la lógica de notificación
		// Pasamos el alert existente para que no se cree duplicado
		err := w.alertService.CreateAlert(ctx, &alert)
		if err != nil {
			log.Printf("Error processing alert %d: %v", alert.ID, err)
		}
	}
}