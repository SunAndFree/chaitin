package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"chaitin-job/work-manager/models"

	"github.com/xuri/excelize/v2"
)

// ExportService handles exporting work arrangements
type ExportService struct {
	svc *WorkArrangementService
}

// NewExportService creates a new export service
func NewExportService(svc *WorkArrangementService) *ExportService {
	return &ExportService{svc: svc}
}

// getFilteredData retrieves data based on filters (all if no filters)
func (e *ExportService) getFilteredData(filters models.FilterParams) ([]models.WorkArrangement, error) {
	hasFilters := filters.DateFrom != "" || filters.DateTo != "" || filters.Customer != "" ||
		filters.Project != "" || filters.WorkType != "" || filters.Progress != ""

	if hasFilters {
		return e.svc.Filter(filters)
	}
	return e.svc.GetAll()
}

// ExportToExcel exports work arrangements to an Excel file
func (e *ExportService) ExportToExcel(filePath string, filters models.FilterParams) error {
	data, err := e.getFilteredData(filters)
	if err != nil {
		return fmt.Errorf("failed to get data for export: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "工作安排"
	f.SetSheetName("Sheet1", sheetName)

	// Define headers
	headers := []string{"ID", "日期", "客户名称", "项目名称", "工作类型", "工作地点", "伙伴", "工作内容", "工作耗时(h)", "工作进度", "备注", "创建时间", "更新时间"}

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "FFFFFF",
			Family: "Microsoft YaHei",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D9D9D9", Style: 1},
			{Type: "top", Color: "D9D9D9", Style: 1},
			{Type: "right", Color: "D9D9D9", Style: 1},
			{Type: "bottom", Color: "D9D9D9", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}

	// Create data cell style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   11,
			Family: "Microsoft YaHei",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D9D9D9", Style: 1},
			{Type: "top", Color: "D9D9D9", Style: 1},
			{Type: "right", Color: "D9D9D9", Style: 1},
			{Type: "bottom", Color: "D9D9D9", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create data style: %w", err)
	}

	// Write headers
	for i, h := range headers {
		cell := fmt.Sprintf("%s1", columnIndexToLetters(i))
		f.SetCellValue(sheetName, cell, h)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write data rows
	for rowIdx, record := range data {
		rowNum := rowIdx + 2
		values := []interface{}{
			record.ProjectID, record.Date, record.Customer, record.Project,
			record.WorkType, record.Location, record.Partner,
			record.Content, record.Duration, record.Progress,
			record.Notes, record.CreatedAt, record.UpdatedAt,
		}
		for colIdx, val := range values {
			cell := fmt.Sprintf("%s%d", columnIndexToLetters(colIdx), rowNum)
			f.SetCellValue(sheetName, cell, val)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Auto-fit column widths
	for i := range headers {
		col := columnIndexToLetters(i)
		// Estimate column width based on header
		width := float64(len(headers[i]) * 3)
		if width < 12 {
			width = 12
		}
		if width > 40 {
			width = 40
		}
		f.SetColWidth(sheetName, col, col, width)
	}

	// Set row height for header
	f.SetRowHeight(sheetName, 1, 28)

	// Save file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// ExportToJSON exports work arrangements to a JSON file
func (e *ExportService) ExportToJSON(filePath string, filters models.FilterParams) error {
	data, err := e.getFilteredData(filters)
	if err != nil {
		return fmt.Errorf("failed to get data for export: %w", err)
	}

	if data == nil {
		data = []models.WorkArrangement{}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// ExportToExcelWriter writes the Excel file to an io.Writer (for HTTP download)
func (e *ExportService) ExportToExcelWriter(w io.Writer, filters models.FilterParams) error {
	data, err := e.getFilteredData(filters)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "工作安排"
	f.SetSheetName("Sheet1", sheetName)

	// Same header and data writing as ExportToExcel but simplified (no styles needed for stream)
	headers := []string{"ID", "日期", "客户名称", "项目名称", "工作类型", "工作地点", "伙伴", "工作内容", "工作耗时(h)", "工作进度", "备注", "创建时间", "更新时间"}

	for i, h := range headers {
		cell := fmt.Sprintf("%s1", columnIndexToLetters(i))
		f.SetCellValue(sheetName, cell, h)
	}

	for rowIdx, record := range data {
		rowNum := rowIdx + 2
		values := []interface{}{
			record.ProjectID, record.Date, record.Customer, record.Project,
			record.WorkType, record.Location, record.Partner,
			record.Content, record.Duration, record.Progress,
			record.Notes, record.CreatedAt, record.UpdatedAt,
		}
		for colIdx, val := range values {
			cell := fmt.Sprintf("%s%d", columnIndexToLetters(colIdx), rowNum)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	return f.Write(w)
}

// ExportToJSONWriter writes JSON to an io.Writer (for HTTP download)
func (e *ExportService) ExportToJSONWriter(w io.Writer, filters models.FilterParams) error {
	data, err := e.getFilteredData(filters)
	if err != nil {
		return err
	}

	if data == nil {
		data = []models.WorkArrangement{}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// columnIndexToLetters converts a 0-based column index to Excel column letters
func columnIndexToLetters(col int) string {
	if col < 26 {
		return string(rune('A' + col))
	}
	return string(rune('A'+col/26-1)) + string(rune('A'+col%26))
}
