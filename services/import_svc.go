package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"chaitin-job/work-manager/models"

	"github.com/xuri/excelize/v2"
)

// ImportService handles importing work arrangements from files
type ImportService struct{}

// NewImportService creates a new import service
func NewImportService() *ImportService {
	return &ImportService{}
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Records []models.WorkArrangement
	Created int
	Skipped int
	Errors  []string
}

// ParseFile parses a file and returns work arrangement records
func (s *ImportService) ParseFile(filePath string) (*ImportResult, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".xlsx":
		return s.parseExcel(filePath)
	case ".json":
		return s.parseJSON(filePath)
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s，请使用 .xlsx 或 .json 文件", ext)
	}
}

// parseExcel parses an Excel file
func (s *ImportService) parseExcel(filePath string) (*ImportResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 Excel 文件: %w", err)
	}
	defer f.Close()

	// Get the first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel 文件中没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("无法读取工作表数据: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel 文件中没有数据行（至少需要标题行和数据行）")
	}

	// Parse header row to determine column mapping
	headerRow := rows[0]
	colMap := make(map[string]int)
	for i, h := range headerRow {
		h = strings.TrimSpace(h)
		if h == "工作耗时(h)" {
			h = "工作耗时" // normalize
		}
		colMap[h] = i
	}

	// Check required columns
	requiredHeaders := []string{"日期", "客户名称", "工作类型", "工作地点"}
	for _, h := range requiredHeaders {
		if _, ok := colMap[h]; !ok {
			return nil, fmt.Errorf("Excel 文件缺少必需的列: %s", h)
		}
	}

	result := &ImportResult{}

	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		record := models.WorkArrangement{}

		// Extract fields
		if col, ok := colMap["project_id"]; ok && col < len(row) {
			var pid int64
			fmt.Sscanf(strings.TrimSpace(row[col]), "%d", &pid)
			record.ProjectID = pid
		}
		if col, ok := colMap["日期"]; ok && col < len(row) {
			record.Date = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["客户名称"]; ok && col < len(row) {
			record.Customer = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["项目名称"]; ok && col < len(row) {
			record.Project = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["工作类型"]; ok && col < len(row) {
			record.WorkType = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["工作地点"]; ok && col < len(row) {
			record.Location = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["伙伴"]; ok && col < len(row) {
			partner := strings.TrimSpace(row[col])
			if partner == "" {
				partner = "否"
			}
			record.Partner = partner
		}
		if col, ok := colMap["工作内容"]; ok && col < len(row) {
			record.Content = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["工作耗时"]; ok && col < len(row) {
			var duration float64
			fmt.Sscanf(strings.TrimSpace(row[col]), "%f", &duration)
			record.Duration = duration
		}
		if col, ok := colMap["工作进度"]; ok && col < len(row) {
			progress := strings.TrimSpace(row[col])
			if progress == "" {
				progress = "未开始"
			}
			record.Progress = progress
		}
		if col, ok := colMap["备注"]; ok && col < len(row) {
			record.Notes = strings.TrimSpace(row[col])
		}

		// Validate record
		if errs := validateRecord(record, rowIdx+2); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.Skipped++
			continue
		}

		result.Records = append(result.Records, record)
		result.Created++
	}

	return result, nil
}

// parseJSON parses a JSON file
func (s *ImportService) parseJSON(filePath string) (*ImportResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取 JSON 文件: %w", err)
	}

	var records []models.WorkArrangement
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("JSON 文件中没有数据记录")
	}

	result := &ImportResult{}

	for i, record := range records {
		// Set defaults for empty fields
		if record.Partner == "" {
			record.Partner = "否"
		}
		if record.Progress == "" {
			record.Progress = "未开始"
		}

		if errs := validateRecord(record, i+2); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.Skipped++
			continue
		}

		result.Records = append(result.Records, record)
		result.Created++
	}

	return result, nil
}

// validateRecord validates a single work arrangement record
func validateRecord(record models.WorkArrangement, rowNum int) []string {
	var errors []string

	if record.Date == "" {
		errors = append(errors, fmt.Sprintf("行 %d: 日期不能为空", rowNum))
	}
	if strings.TrimSpace(record.Customer) == "" {
		errors = append(errors, fmt.Sprintf("行 %d: 客户名称不能为空", rowNum))
	}
	if !models.IsValidWorkType(record.WorkType) {
		errors = append(errors, fmt.Sprintf("行 %d: 无效的工作类型 '%s'", rowNum, record.WorkType))
	}
	if !models.IsValidLocation(record.Location) {
		errors = append(errors, fmt.Sprintf("行 %d: 无效的工作地点 '%s'", rowNum, record.Location))
	}
	if !models.IsValidPartner(record.Partner) {
		errors = append(errors, fmt.Sprintf("行 %d: 无效的伙伴值 '%s'", rowNum, record.Partner))
	}
	if !models.IsValidProgress(record.Progress) {
		errors = append(errors, fmt.Sprintf("行 %d: 无效的工作进度 '%s'", rowNum, record.Progress))
	}

	return errors
}

// ParseReader parses uploaded file content (from multipart upload) into work arrangement records.
func (s *ImportService) ParseReader(reader io.Reader, filename string) (*ImportResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".xlsx":
		return s.parseExcelReader(reader)
	case ".json":
		return s.parseJSONReader(reader)
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s，请使用 .xlsx 或 .json 文件", ext)
	}
}

func (s *ImportService) parseExcelReader(reader io.Reader) (*ImportResult, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("无法读取 Excel 文件: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel 文件中没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("无法读取工作表数据: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel 文件中没有数据行（至少需要标题行和数据行）")
	}

	return parseRows(rows)
}

func (s *ImportService) parseJSONReader(reader io.Reader) (*ImportResult, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("无法读取 JSON 数据: %w", err)
	}

	var records []models.WorkArrangement
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("JSON 文件中没有数据记录")
	}

	result := &ImportResult{}
	for i, record := range records {
		if record.Partner == "" {
			record.Partner = "否"
		}
		if record.Progress == "" {
			record.Progress = "未开始"
		}
		if errs := validateRecord(record, i+2); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.Skipped++
			continue
		}
		result.Records = append(result.Records, record)
		result.Created++
	}
	return result, nil
}

// parseRows parses Excel rows using header-based column mapping (shared by file and reader paths).
func parseRows(rows [][]string) (*ImportResult, error) {
	headerRow := rows[0]
	colMap := make(map[string]int)
	for i, h := range headerRow {
		h = strings.TrimSpace(h)
		if h == "ID" {
			h = "project_id"
		}
		if h == "工作耗时(h)" {
			h = "工作耗时"
		}
		colMap[h] = i
	}

	requiredHeaders := []string{"日期", "客户名称", "工作类型", "工作地点"}
	for _, h := range requiredHeaders {
		if _, ok := colMap[h]; !ok {
			return nil, fmt.Errorf("Excel 文件缺少必需的列: %s", h)
		}
	}

	result := &ImportResult{}
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		record := models.WorkArrangement{}

		if col, ok := colMap["project_id"]; ok && col < len(row) {
			var pid int64
			fmt.Sscanf(strings.TrimSpace(row[col]), "%d", &pid)
			record.ProjectID = pid
		}
		if col, ok := colMap["日期"]; ok && col < len(row) {
			record.Date = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["客户名称"]; ok && col < len(row) {
			record.Customer = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["项目名称"]; ok && col < len(row) {
			record.Project = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["工作类型"]; ok && col < len(row) {
			record.WorkType = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["工作地点"]; ok && col < len(row) {
			record.Location = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["伙伴"]; ok && col < len(row) {
			partner := strings.TrimSpace(row[col])
			if partner == "" {
				partner = "否"
			}
			record.Partner = partner
		}
		if col, ok := colMap["工作内容"]; ok && col < len(row) {
			record.Content = strings.TrimSpace(row[col])
		}
		if col, ok := colMap["工作耗时"]; ok && col < len(row) {
			var duration float64
			fmt.Sscanf(strings.TrimSpace(row[col]), "%f", &duration)
			record.Duration = duration
		}
		if col, ok := colMap["工作进度"]; ok && col < len(row) {
			progress := strings.TrimSpace(row[col])
			if progress == "" {
				progress = "未开始"
			}
			record.Progress = progress
		}
		if col, ok := colMap["备注"]; ok && col < len(row) {
			record.Notes = strings.TrimSpace(row[col])
		}

		if errs := validateRecord(record, rowIdx+2); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.Skipped++
			continue
		}
		result.Records = append(result.Records, record)
		result.Created++
	}
	return result, nil
}
