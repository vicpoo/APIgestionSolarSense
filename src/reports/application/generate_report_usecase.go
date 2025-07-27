// src/reports/application/generate_report_usecase.go
package application

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
        return 0, fmt.Errorf("no se pudo obtener el ID del usuario: %v. Nota: Asegúrate que el email existe en la tabla users", err)
    }
    return userID, nil
}

func (uc *GenerateReportUseCase) GeneratePDFReport(ctx context.Context, req domain.GenerateReportRequest) (*domain.Report, error) {
    // 1. Obtener ID del usuario
    userID, err := uc.GetUserIDByEmail(ctx, req.RequestedByEmail)
    if err != nil {
        return nil, fmt.Errorf("error al obtener ID de usuario: %v", err)
    }

    // 2. Obtener lecturas de la fecha solicitada (últimos 50)
    readings, err := uc.repo.GetSensorReadingsByDate(ctx, req.Date)
    if err != nil {
        return nil, fmt.Errorf("error al obtener lecturas: %v", err)
    }

    // 3. Generar PDF con mejor formato
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    
    // Encabezado
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(0, 10, "Reporte de Sensores - "+req.Date)
    pdf.Ln(15)
    
    // Información del solicitante
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(0, 10, "Solicitado por: "+req.RequestedByEmail)
    pdf.Ln(10)
    pdf.Cell(0, 10, fmt.Sprintf("Total de registros: %d", len(readings)))
    pdf.Ln(15)

    // Configurar tabla
    header := []string{"Hora", "Sensor ID", "Temp (°C)", "Hum (%)", "Volt (V)", "Corr (A)", "Pot (W)"}
    colWidths := []float64{30, 20, 20, 20, 20, 20, 20}
    
    // Estilo de tabla
    pdf.SetFont("Arial", "B", 10)
    pdf.SetFillColor(200, 200, 200)
    
    // Encabezados de tabla
    for i, str := range header {
        pdf.CellFormat(colWidths[i], 7, str, "1", 0, "C", true, 0, "")
    }
    pdf.Ln(-1)
    
    // Datos de la tabla
    pdf.SetFont("Arial", "", 8)
    pdf.SetFillColor(255, 255, 255)
    
    for _, r := range readings {
        // Formatear valores NULL
        temp := formatNullableFloat(r.Temperature, "%.1f")
        hum := formatNullableFloat(r.Humidity, "%.1f")
        volt := formatNullableFloat(r.Voltage, "%.2f")
        curr := formatNullableFloat(r.Current, "%.2f")
        pot := formatNullableFloat(r.Potencia, "%.2f")
        
        // Fila de datos
        pdf.CellFormat(colWidths[0], 6, r.RecordedAt.Format("15:04:05"), "1", 0, "C", false, 0, "")
        pdf.CellFormat(colWidths[1], 6, fmt.Sprintf("%d", r.SensorID), "1", 0, "C", false, 0, "")
        pdf.CellFormat(colWidths[2], 6, temp, "1", 0, "C", false, 0, "")
        pdf.CellFormat(colWidths[3], 6, hum, "1", 0, "C", false, 0, "")
        pdf.CellFormat(colWidths[4], 6, volt, "1", 0, "C", false, 0, "")
        pdf.CellFormat(colWidths[5], 6, curr, "1", 0, "C", false, 0, "")
        pdf.CellFormat(colWidths[6], 6, pot, "1", 0, "C", false, 0, "")
        pdf.Ln(-1)
    }

    // 6. Guardar PDF en el sistema de archivos
    fileName := fmt.Sprintf("reporte_%s_%s.pdf", req.Date, time.Now().Format("20060102_150405"))
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

    return report, nil
}

// Función auxiliar para formatear valores float que pueden ser NULL
func formatNullableFloat(f *float64, format string) string {
    if f == nil {
        return "N/A"
    }
    return fmt.Sprintf(format, *f)
}