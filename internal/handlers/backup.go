package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type BackupHandler struct {
	backupService *services.BackupService
}

func NewBackupHandler(s *services.BackupService) *BackupHandler {
	return &BackupHandler{backupService: s}
}

func (h *BackupHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId query parameter")
		return
	}
	list, err := h.backupService.ListConfigsByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *BackupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cfg models.BackupConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.backupService.CreateConfig(r.Context(), &cfg); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, cfg)
}

func (h *BackupHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	cfg, err := h.backupService.GetConfig(r.Context(), id)
	if err != nil || cfg == nil {
		WriteError(w, http.StatusNotFound, "backup config not found")
		return
	}
	WriteJSON(w, http.StatusOK, cfg)
}

func (h *BackupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projectID := r.URL.Query().Get("projectId")
	if id == "" || projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing id or projectId")
		return
	}
	if err := h.backupService.DeleteConfig(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BackupHandler) Trigger(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	rec, err := h.backupService.TriggerBackup(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, rec)
}

func (h *BackupHandler) ListRecords(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	recs, err := h.backupService.ListRecordsByConfig(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, recs)
}

func (h *BackupHandler) ListS3Destinations(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId query parameter")
		return
	}
	list, err := h.backupService.ListS3Destinations(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *BackupHandler) CreateS3Destination(w http.ResponseWriter, r *http.Request) {
	var dest models.S3Destination
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.backupService.CreateS3Destination(r.Context(), &dest); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, dest)
}

func (h *BackupHandler) DeleteS3Destination(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projectID := r.URL.Query().Get("projectId")
	if id == "" || projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing id or projectId")
		return
	}
	if err := h.backupService.DeleteS3Destination(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
