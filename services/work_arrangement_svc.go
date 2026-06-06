package services

import (
	"database/sql"
	"fmt"
	"strings"

	"chaitin-job/work-manager/db"
	"chaitin-job/work-manager/models"
)

type WorkArrangementService struct{}

func NewWorkArrangementService() *WorkArrangementService {
	return &WorkArrangementService{}
}

func (s *WorkArrangementService) GetAll() ([]models.WorkArrangement, error) {
	rows, err := db.DB.Query(`
		SELECT id, project_id, date, customer, project, work_type, location, partner,
		       content, duration, progress, notes, created_at, updated_at
		FROM work_arrangements ORDER BY date DESC, updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	return scanWorkArrangements(rows)
}

func (s *WorkArrangementService) GetByID(id int64) (*models.WorkArrangement, error) {
	row := db.DB.QueryRow(`
		SELECT id, project_id, date, customer, project, work_type, location, partner,
		       content, duration, progress, notes, created_at, updated_at
		FROM work_arrangements WHERE id = ?
	`, id)
	return scanWorkArrangement(row)
}

func (s *WorkArrangementService) Create(wa models.WorkArrangement) (*models.WorkArrangement, error) {
	if !models.IsValidWorkType(wa.WorkType) {
		return nil, fmt.Errorf("invalid work_type: %s", wa.WorkType)
	}
	if !models.IsValidLocation(wa.Location) {
		return nil, fmt.Errorf("invalid location: %s", wa.Location)
	}
	if !models.IsValidPartner(wa.Partner) {
		return nil, fmt.Errorf("invalid partner: %s", wa.Partner)
	}
	if !models.IsValidProgress(wa.Progress) {
		return nil, fmt.Errorf("invalid progress: %s", wa.Progress)
	}

	result, err := db.DB.Exec(`
		INSERT INTO work_arrangements (project_id, date, customer, project, work_type, location, partner, content, duration, progress, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, wa.ProjectID, wa.Date, wa.Customer, wa.Project, wa.WorkType, wa.Location, wa.Partner, wa.Content, wa.Duration, wa.Progress, wa.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get id: %w", err)
	}
	return s.GetByID(id)
}

func (s *WorkArrangementService) Update(wa models.WorkArrangement) (*models.WorkArrangement, error) {
	if !models.IsValidWorkType(wa.WorkType) {
		return nil, fmt.Errorf("invalid work_type: %s", wa.WorkType)
	}
	if !models.IsValidLocation(wa.Location) {
		return nil, fmt.Errorf("invalid location: %s", wa.Location)
	}
	if !models.IsValidPartner(wa.Partner) {
		return nil, fmt.Errorf("invalid partner: %s", wa.Partner)
	}
	if !models.IsValidProgress(wa.Progress) {
		return nil, fmt.Errorf("invalid progress: %s", wa.Progress)
	}

	existing, err := s.GetByID(wa.ID)
	if err != nil {
		return nil, fmt.Errorf("not found: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("work arrangement not found: id=%d", wa.ID)
	}

	createdAt := wa.CreatedAt
	if createdAt == "" {
		createdAt = existing.CreatedAt
	}

	_, err = db.DB.Exec(`
		UPDATE work_arrangements
		SET project_id = ?, date = ?, customer = ?, project = ?, work_type = ?, location = ?,
		    partner = ?, content = ?, duration = ?, progress = ?, notes = ?,
		    created_at = ?, updated_at = datetime('now', 'localtime')
		WHERE id = ?
	`, wa.ProjectID, wa.Date, wa.Customer, wa.Project, wa.WorkType, wa.Location,
		wa.Partner, wa.Content, wa.Duration, wa.Progress, wa.Notes, createdAt, wa.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update: %w", err)
	}
	return s.GetByID(wa.ID)
}

func (s *WorkArrangementService) Delete(id int64) error {
	result, err := db.DB.Exec(`DELETE FROM work_arrangements WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("not found: id=%d", id)
	}
	return nil
}

func (s *WorkArrangementService) Filter(filters models.FilterParams) ([]models.WorkArrangement, error) {
	query := `SELECT id, project_id, date, customer, project, work_type, location, partner,
	                 content, duration, progress, notes, created_at, updated_at
	          FROM work_arrangements WHERE 1=1`
	args := []interface{}{}

	if filters.DateFrom != "" {
		query += " AND date >= ?"
		args = append(args, filters.DateFrom)
	}
	if filters.DateTo != "" {
		query += " AND date <= ?"
		args = append(args, filters.DateTo)
	}
	if filters.Customer != "" {
		query += " AND customer LIKE ?"
		args = append(args, "%"+filters.Customer+"%")
	}
	if filters.Project != "" {
		query += " AND project LIKE ?"
		args = append(args, "%"+filters.Project+"%")
	}
	if filters.WorkType != "" {
		query += " AND work_type = ?"
		args = append(args, filters.WorkType)
	}
	if filters.Progress != "" {
		query += " AND progress LIKE ?"
		args = append(args, "%"+filters.Progress+"%")
	}

	query += " ORDER BY date DESC, updated_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("filter failed: %w", err)
	}
	defer rows.Close()
	return scanWorkArrangements(rows)
}

func (s *WorkArrangementService) BulkCreate(records []models.WorkArrangement) (int, int, []string, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return 0, 0, nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO work_arrangements (project_id, date, customer, project, work_type, location, partner, content, duration, progress, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	created, skipped := 0, 0
	var errors []string

	for i, record := range records {
		if !models.IsValidWorkType(record.WorkType) {
			errors = append(errors, fmt.Sprintf("行 %d: 无效类型 '%s'", i+2, record.WorkType))
			skipped++
			continue
		}
		if !models.IsValidLocation(record.Location) {
			errors = append(errors, fmt.Sprintf("行 %d: 无效地点 '%s'", i+2, record.Location))
			skipped++
			continue
		}
		if !models.IsValidPartner(record.Partner) {
			errors = append(errors, fmt.Sprintf("行 %d: 无效伙伴 '%s'", i+2, record.Partner))
			skipped++
			continue
		}
		if !models.IsValidProgress(record.Progress) {
			errors = append(errors, fmt.Sprintf("行 %d: 无效进度 '%s'", i+2, record.Progress))
			skipped++
			continue
		}
		if record.Date == "" {
			errors = append(errors, fmt.Sprintf("行 %d: 日期为空", i+2))
			skipped++
			continue
		}
		if strings.TrimSpace(record.Customer) == "" {
			errors = append(errors, fmt.Sprintf("行 %d: 客户为空", i+2))
			skipped++
			continue
		}
		_, err := stmt.Exec(
			record.ProjectID, record.Date, record.Customer, record.Project, record.WorkType,
			record.Location, record.Partner, record.Content, record.Duration,
			record.Progress, record.Notes,
		)
		if err != nil {
			errors = append(errors, fmt.Sprintf("行 %d: %v", i+2, err))
			skipped++
			continue
		}
		created++
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, nil, fmt.Errorf("commit: %w", err)
	}
	return created, skipped, errors, nil
}

func (s *WorkArrangementService) GetDistinctCustomers() ([]string, error) {
	rows, err := db.DB.Query(`SELECT DISTINCT customer FROM work_arrangements WHERE customer != '' ORDER BY customer`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var c string
		rows.Scan(&c)
		res = append(res, c)
	}
	return res, rows.Err()
}

func (s *WorkArrangementService) GetDistinctProjects() ([]string, error) {
	rows, err := db.DB.Query(`SELECT DISTINCT project FROM work_arrangements WHERE project != '' ORDER BY project`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		res = append(res, p)
	}
	return res, rows.Err()
}

func (s *WorkArrangementService) GenerateCopyText(id int64) (string, error) {
	wa, err := s.GetByID(id)
	if err != nil {
		return "", err
	}
	if wa == nil {
		return "", fmt.Errorf("not found: %d", id)
	}
	return wa.GenerateCopyText(), nil
}

func scanWorkArrangements(rows *sql.Rows) ([]models.WorkArrangement, error) {
	var result []models.WorkArrangement
	for rows.Next() {
		var wa models.WorkArrangement
		err := rows.Scan(&wa.ID, &wa.ProjectID, &wa.Date, &wa.Customer, &wa.Project,
			&wa.WorkType, &wa.Location, &wa.Partner, &wa.Content, &wa.Duration,
			&wa.Progress, &wa.Notes, &wa.CreatedAt, &wa.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, wa)
	}
	if result == nil {
		result = []models.WorkArrangement{}
	}
	return result, rows.Err()
}

func scanWorkArrangement(row *sql.Row) (*models.WorkArrangement, error) {
	var wa models.WorkArrangement
	err := row.Scan(&wa.ID, &wa.ProjectID, &wa.Date, &wa.Customer, &wa.Project,
		&wa.WorkType, &wa.Location, &wa.Partner, &wa.Content, &wa.Duration,
		&wa.Progress, &wa.Notes, &wa.CreatedAt, &wa.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	return &wa, nil
}
