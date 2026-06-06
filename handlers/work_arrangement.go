package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chaitin-job/work-manager/models"
	"chaitin-job/work-manager/services"
)

type WorkArrangementHandler struct {
	svc *services.WorkArrangementService
}

func NewWorkArrangementHandler(svc *services.WorkArrangementService) *WorkArrangementHandler {
	return &WorkArrangementHandler{svc: svc}
}

// List returns all work arrangements, optionally filtered
func (h *WorkArrangementHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := models.FilterParams{
		DateFrom: q.Get("date_from"),
		DateTo:   q.Get("date_to"),
		Customer: q.Get("customer"),
		Project:  q.Get("project"),
		WorkType: q.Get("work_type"),
		Progress: q.Get("progress"),
	}

	data, err := h.svc.Filter(filters)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// Get returns a single work arrangement by ID
func (h *WorkArrangementHandler) Get(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	record, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if record == nil {
		writeError(w, http.StatusNotFound, "record not found")
		return
	}
	writeJSON(w, http.StatusOK, record)
}

// Create inserts a new work arrangement
func (h *WorkArrangementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var record models.WorkArrangement
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	result, err := h.svc.Create(record)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// Update modifies an existing work arrangement
func (h *WorkArrangementHandler) Update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var record models.WorkArrangement
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	record.ID = id

	result, err := h.svc.Update(record)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Delete removes a work arrangement by ID
func (h *WorkArrangementHandler) Delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// BulkCreate inserts multiple work arrangements
func (h *WorkArrangementHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var records []models.WorkArrangement
	if err := json.NewDecoder(r.Body).Decode(&records); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	created, skipped, errors, err := h.svc.BulkCreate(records)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"created": created,
		"skipped": skipped,
		"errors":  errors,
	})
}

// Customers returns distinct customer names
func (h *WorkArrangementHandler) Customers(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetDistinctCustomers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// Projects returns distinct project names
func (h *WorkArrangementHandler) Projects(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetDistinctProjects()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// writeJSON helper
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError helper
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
