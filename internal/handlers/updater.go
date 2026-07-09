package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/services"
)

type UpdaterHandler struct {
	updaterService *services.UpdaterService
}

func NewUpdaterHandler(s *services.UpdaterService) *UpdaterHandler {
	return &UpdaterHandler{updaterService: s}
}

func (h *UpdaterHandler) GetUpdateStatus(c echo.Context) error {
	if h.updaterService == nil {
		WriteError(w, http.StatusInternalServerError, "updater service not initialized")
		return nil
	}
	status := h.updaterService.GetStatus()
	WriteJSON(w, http.StatusOK, status)
}

func (h *UpdaterHandler) CheckUpdate(c echo.Context) error {
	if h.updaterService == nil {
		WriteError(w, http.StatusInternalServerError, "updater service not initialized")
		return nil
	}
	if _, err := h.updaterService.CheckForUpdates(r.Context()); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	status := h.updaterService.GetStatus()
	WriteJSON(w, http.StatusOK, status)
}

func (h *UpdaterHandler) DeployUpdate(c echo.Context) error {
	if h.updaterService == nil {
		WriteError(w, http.StatusInternalServerError, "updater service not initialized")
		return nil
	}
	if err := h.updaterService.DeployUpdate(r.Context()); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusAccepted, map[string]string{
		"message": "update deployment triggered",
	})
}
