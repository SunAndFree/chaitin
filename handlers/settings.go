package handlers

import (
	"encoding/json"
	"net/http"

	"chaitin-job/work-manager/platform"
)

type SettingsHandler struct{}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

func (h *SettingsHandler) GetAutoStart(w http.ResponseWriter, r *http.Request) {
	enabled, err := platform.IsAutoStartEnabled()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"enabled": enabled})
}

func (h *SettingsHandler) SetAutoStart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	var err error
	if body.Enabled {
		err = platform.EnableAutoStart()
	} else {
		err = platform.DisableAutoStart()
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"enabled": body.Enabled})
}
