package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// runMigrations creates or upgrades the database schema.
func runMigrations(database *sql.DB) error {
	// Check if project_id column exists
	hasProjectID := false
	rows, err := database.Query(`PRAGMA table_info(work_arrangements)`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cid int
			var name, typ string
			var notNull int
			var dflt, pk sql.NullString
			rows.Scan(&cid, &name, &typ, &notNull, &dflt, &pk)
			if name == "project_id" {
				hasProjectID = true
			}
		}
	}

	if !hasProjectID {
		// Add project_id column for existing databases
		database.Exec(`ALTER TABLE work_arrangements ADD COLUMN project_id INTEGER NOT NULL DEFAULT 0`)
	}

	// Fresh install - create table if not exists
	creates := []string{
		`CREATE TABLE IF NOT EXISTS work_arrangements (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id  INTEGER NOT NULL DEFAULT 0,
			date        TEXT NOT NULL,
			customer    TEXT NOT NULL DEFAULT '',
			project     TEXT NOT NULL DEFAULT '',
			work_type   TEXT NOT NULL CHECK(work_type IN ('测试','交付','售后')),
			location    TEXT NOT NULL CHECK(location IN ('远程','现场')),
			partner     TEXT NOT NULL CHECK(partner IN ('是','否')),
			content     TEXT NOT NULL DEFAULT '',
			duration    REAL NOT NULL DEFAULT 0,
			progress    TEXT NOT NULL CHECK(progress IN ('未开始','进行中','已完成','已暂停','已取消')),
			notes       TEXT NOT NULL DEFAULT '',
			created_at  TEXT NOT NULL DEFAULT (datetime('now','localtime')),
			updated_at  TEXT NOT NULL DEFAULT (datetime('now','localtime'))
		)`,
	}
	for _, m := range creates {
		if _, err := database.Exec(m); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("migration failed: %w\nSQL: %s", err, m)
			}
		}
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_date ON work_arrangements(date)`,
		`CREATE INDEX IF NOT EXISTS idx_customer ON work_arrangements(customer)`,
		`CREATE INDEX IF NOT EXISTS idx_project ON work_arrangements(project)`,
		`CREATE INDEX IF NOT EXISTS idx_work_type ON work_arrangements(work_type)`,
		`CREATE INDEX IF NOT EXISTS idx_progress ON work_arrangements(progress)`,
	}
	for _, m := range indexes {
		if _, err := database.Exec(m); err != nil {
			return fmt.Errorf("index creation failed: %w\nSQL: %s", err, m)
		}
	}

	return nil
}
