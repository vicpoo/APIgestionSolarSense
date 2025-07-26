// src/reports/application/generate_report_usecase.go
package application

import (
	"context"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/vicpoo/apigestion-solar-go/src/reports/domain"
)

type GenerateReportUseCase struct {
    repo domain.ReportRepository
}

func NewGenerateReportUseCase(repo domain.ReportRepository) *GenerateReportUseCase {
    return &GenerateReportUseCase{repo: repo}
}

func (uc *GenerateReportUseCase) GeneratePDFReport(ctx context.Context, req domain.GenerateReportRequest) (*domain.Report, error) {
    // 1. Obtener lecturas de la fecha solicitada
    readings, err := uc.repo.GetSensorReadingsByDate(ctx, req.Date)
    if err != nil {
        return nil, err
    }

    // 2. Generar PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Reporte de Sensores - "+req.Date)
    pdf.Ln(20)

    // 3. Agregar datos del usuario solicitante
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, "Solicitado por: "+req.RequestedByEmail)
    pdf.Ln(10)

    // 4. Tabla con lecturas (ejemplo simplificado)
    pdf.Cell(40, 10, "Sensor ID | Temperatura | Humedad | Voltaje")
    pdf.Ln(10)
for _, r := range readings {
    // Convertir valores numéricos a string (manejando punteros nulos)
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

    // Construir línea del PDF
    pdf.Cell(40, 10, fmt.Sprintf(
        "%d | %s | %s | %s",
        r.SensorID, // SensorID es int, no necesita conversión
        temp,
        hum,
        volt,
    ))
    pdf.Ln(10)
}

    // 5. Guardar PDF en el sistema de archivos
    fileName := "reporte_" + req.Date + ".pdf"
    storagePath := "/storage/reports/" + fileName
    err = pdf.OutputFileAndClose(storagePath)
    if err != nil {
        return nil, err
    }

    // 6. Guardar metadata en la BD
    report := &domain.Report{
        UserID:      1, // ID del admin (dueño de los sensores)
        FileName:    fileName,
        StoragePath: storagePath,
        Format:      "PDF",
    }
    err = uc.repo.Create(ctx, report)
    return report, err
}