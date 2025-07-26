// src/reports/application/generate_report_usecase.go
package application

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type GenerateReportUseCase struct {
    repo domain.ReportRepository
    db   *sql.DB // Conexión a la base de datos para consultas de usuario
}

// Modificado para recibir la conexión a la base de datos
func NewGenerateReportUseCase(repo domain.ReportRepository, db *sql.DB) *GenerateReportUseCase {
    return &GenerateReportUseCase{repo: repo, db: db}
}

func (uc *GenerateReportUseCase) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
    query := `SELECT id FROM users WHERE email = ?`
    var userID int
    err := uc.db.QueryRowContext(ctx, query, email).Scan(&userID)
    if err != nil {
        return 0, fmt.Errorf("no se pudo obtener el ID del usuario: %v", err)
    }
    return userID, nil
}

func (uc *GenerateReportUseCase) GeneratePDFReport(ctx context.Context, req domain.GenerateReportRequest) (*domain.Report, error) {
    // 1. Obtener ID del usuario
    userID, err := uc.GetUserIDByEmail(ctx, req.RequestedByEmail)
    if err != nil {
        return nil, fmt.Errorf("error al obtener ID de usuario: %v", err)
    }

    // 2. Obtener lecturas de la fecha solicitada
    readings, err := uc.repo.GetSensorReadingsByDate(ctx, req.Date)
    if err != nil {
        return nil, fmt.Errorf("error al obtener lecturas: %v", err)
    }

    // 3. Generar PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Reporte de Sensores - "+req.Date)
    pdf.Ln(20)

    // 4. Agregar datos del usuario solicitante
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, "Solicitado por: "+req.RequestedByEmail)
    pdf.Ln(10)

    // 5. Tabla con lecturas
    pdf.Cell(40, 10, "Sensor ID | Temperatura | Humedad | Voltaje")
    pdf.Ln(10)
    
    for _, r := range readings {
        temp := "N/A"
        if r.Temperature != nil {
            temp = fmt.Sprintf("%.2f°C", *r.Temperature)
        }
        
        hum := "N/A"
        if r.Humidity != nil {
            hum = fmt.Sprintf("%.2f%%", *r.Humidity)
        }
        
        volt := "N/A"
        if r.Voltage != nil {
            volt = fmt.Sprintf("%.2fV", *r.Voltage)
        }

        pdf.Cell(40, 10, fmt.Sprintf(
            "%d | %s | %s | %s",
            r.SensorID,
            temp,
            hum,
            volt,
        ))
        pdf.Ln(10)
    }

    // 6. Guardar PDF en el sistema de archivos
    fileName := "reporte_" + req.Date + ".pdf"
    storagePath := "/home/ubuntu/APIgestionSolarSense/storage/reports/" + fileName
    err = pdf.OutputFileAndClose(storagePath)
    if err != nil {
        return nil, fmt.Errorf("error al guardar PDF: %v", err)
    }

    // 7. Guardar metadata en la BD
    report := &domain.Report{
        UserID:      userID,
        FileName:    fileName,
        StoragePath: storagePath,
        Format:      "PDF",
    }
    err = uc.repo.Create(ctx, report)
    if err != nil {
        return nil, fmt.Errorf("error al guardar reporte en BD: %v", err)
    }
    if report.ID == 0 {
        return nil, fmt.Errorf("no se pudo obtener el ID del reporte generado")
    }

    return report, nil
}