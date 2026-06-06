package handlers

import (
	"encoding/json"
	"net/http"

	"chaitin-job/work-manager/models"
	"chaitin-job/work-manager/services"
)

type ImportHandler struct {
	importSvc *services.ImportService
	svc       *services.WorkArrangementService
}

func NewImportHandler(importSvc *services.ImportService, svc *services.WorkArrangementService) *ImportHandler {
	return &ImportHandler{importSvc: importSvc, svc: svc}
}

// Parse accepts a multipart file upload and returns parsed records for preview
func (h *ImportHandler) Parse(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "文件过大或格式错误: "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "请选择要上传的文件: "+err.Error())
		return
	}
	defer file.Close()

	result, err := h.importSvc.ParseReader(file, header.Filename)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Confirm accepts parsed records and saves them to the database
func (h *ImportHandler) Confirm(w http.ResponseWriter, r *http.Request) {
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
