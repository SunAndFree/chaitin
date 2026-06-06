package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"chaitin-job/work-manager/db"
	"chaitin-job/work-manager/handlers"
	"chaitin-job/work-manager/services"
)

var (
	svc       *services.WorkArrangementService
	exportSvc *services.ExportService
	importSvc *services.ImportService
)

func initApp() error {
	if err := db.InitDB(); err != nil {
		return err
	}
	svc = services.NewWorkArrangementService()
	exportSvc = services.NewExportService(svc)
	importSvc = services.NewImportService()
	return nil
}

func closeDB() {
	db.CloseDB()
}

func newServer(frontendFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	// Register API handlers
	waHandler := handlers.NewWorkArrangementHandler(svc)
	exportHandler := handlers.NewExportHandler(exportSvc, svc)
	importHandler := handlers.NewImportHandler(importSvc, svc)
	settingsHandler := handlers.NewSettingsHandler()

	// Work arrangements CRUD
	mux.HandleFunc("/api/work-arrangements", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		switch r.Method {
		case http.MethodGet:
			waHandler.List(w, r)
		case http.MethodPost:
			waHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/work-arrangements/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api/work-arrangements/")

		if path == "bulk" && r.Method == http.MethodPost {
			waHandler.BulkCreate(w, r)
			return
		}

		// /api/work-arrangements/{id}
		switch r.Method {
		case http.MethodGet:
			waHandler.Get(w, r, path)
		case http.MethodPut:
			waHandler.Update(w, r, path)
		case http.MethodDelete:
			waHandler.Delete(w, r, path)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Export
	mux.HandleFunc("/api/export/excel", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		exportHandler.ExportExcel(w, r)
	})
	mux.HandleFunc("/api/export/json", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		exportHandler.ExportJSON(w, r)
	})

	// Import
	mux.HandleFunc("/api/import/parse", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		importHandler.Parse(w, r)
	})
	mux.HandleFunc("/api/import/confirm", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		importHandler.Confirm(w, r)
	})

	// Reference data
	mux.HandleFunc("/api/reference/customers", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		waHandler.Customers(w, r)
	})
	mux.HandleFunc("/api/reference/projects", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		waHandler.Projects(w, r)
	})

	// Settings
	mux.HandleFunc("/api/settings/autostart", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		switch r.Method {
		case http.MethodGet:
			settingsHandler.GetAutoStart(w, r)
		case http.MethodPut:
			settingsHandler.SetAutoStart(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve embedded frontend static files
	mux.Handle("/", http.FileServer(http.FS(frontendFS)))

	return mux
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
