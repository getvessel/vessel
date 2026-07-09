package handlers

import (
	"github.com/labstack/echo/v4"

	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type StorageHandler struct {
	storageService *services.StorageService
}

func NewStorageHandler(s *services.StorageService) *StorageHandler {
	return &StorageHandler{storageService: s}
}

func (h *StorageHandler) ListStorage(c echo.Context) error {
	storages, err := h.storageService.ListStorage(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	if storages == nil {
		storages = []*models.Storage{}
	}
	WriteJSON(w, http.StatusOK, storages)
}

func (h *StorageHandler) CreateStorage(c echo.Context) error {
	var st models.Storage
	if err := c.Bind(&st); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	created, err := h.storageService.CreateStorageWithDefaults(r.Context(), &st)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *StorageHandler) GetStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return nil
	}
	st, err := h.storageService.GetStorage(r.Context(), id)
	if err != nil || st == nil {
		WriteError(w, http.StatusNotFound, "storage record not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, st)
}

func (h *StorageHandler) DeleteStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return nil
	}
	if err := h.storageService.DeleteStorage(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *StorageHandler) StartStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return nil
	}
	st, err := h.storageService.StartStorage(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, st)
}

func (h *StorageHandler) StopStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return nil
	}
	if err := h.storageService.StopStorage(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
