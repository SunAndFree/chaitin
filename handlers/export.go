package handlers

import (
	"net/http"

	"chaitin-job/work-manager/models"
	"chaitin-job/work-manager/services"
)

type ExportHandler struct {
	exportSvc *services.ExportService
	svc       *services.WorkArrangementService
}

func NewExportHandler(exportSvc *services.ExportService, svc *services.WorkArrangementService) *ExportHandler {
	return &ExportHandler{exportSvc: exportSvc, svc: svc}
}

func (h *ExportHandler) ExportExcel(w http.ResponseWriter, r *http.Request) {
	filters := parseFilters(r)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=工作安排.xlsx")
	if err := h.exportSvc.ExportToExcelWriter(w, filters); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *ExportHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	filters := parseFilters(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=工作安排.json")
	if err := h.exportSvc.ExportToJSONWriter(w, filters); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func parseFilters(r *http.Request) models.FilterParams {
	q := r.URL.Query()
	return models.FilterParams{
		DateFrom: q.Get("date_from"),
		DateTo:   q.Get("date_to"),
		Customer: q.Get("customer"),
		Project:  q.Get("project"),
		WorkType: q.Get("work_type"),
		Progress: q.Get("progress"),
	}
}
