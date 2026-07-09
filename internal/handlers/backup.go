package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *BackupHandler) List(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId query parameter")
		return nil
	}
	list, err := h.backupService.ListConfigsByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *BackupHandler) Create(c echo.Context) error {
	var cfg models.BackupConfig
	if err := c.Bind(&cfg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.backupService.CreateConfig(r.Context(), &cfg); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, cfg)
}

func (h *BackupHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	cfg, err := h.backupService.GetConfig(r.Context(), id)
	if err != nil || cfg == nil {
		WriteError(w, http.StatusNotFound, "backup config not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, cfg)
}

func (h *BackupHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	projectID := c.QueryParam("projectId")
	if id == "" || projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing id or projectId")
		return nil
	}
	if err := h.backupService.DeleteConfig(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BackupHandler) Trigger(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	rec, err := h.backupService.TriggerBackup(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, rec)
}

func (h *BackupHandler) ListRecords(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	recs, err := h.backupService.ListRecordsByConfig(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, recs)
}

func (h *BackupHandler) ListS3Destinations(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId query parameter")
		return nil
	}
	list, err := h.backupService.ListS3Destinations(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *BackupHandler) CreateS3Destination(c echo.Context) error {
	var dest models.S3Destination
	if err := c.Bind(&dest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.backupService.CreateS3Destination(r.Context(), &dest); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, dest)
}

func (h *BackupHandler) DeleteS3Destination(c echo.Context) error {
	id := c.Param("id")
	projectID := c.QueryParam("projectId")
	if id == "" || projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing id or projectId")
		return nil
	}
	if err := h.backupService.DeleteS3Destination(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}
