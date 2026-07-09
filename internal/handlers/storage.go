package handlers

import (
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

func (h *StorageHandler) ListStorage(w http.ResponseWriter, r *http.Request) {
	storages, err := h.storageService.ListStorage(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if storages == nil {
		storages = []*models.Storage{}
	}
	WriteJSON(w, http.StatusOK, storages)
}

func (h *StorageHandler) CreateStorage(w http.ResponseWriter, r *http.Request) {
	var st models.Storage
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid storage configuration payload")
		return
	}
	created, err := h.storageService.CreateStorageWithDefaults(r.Context(), &st)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *StorageHandler) GetStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}
	st, err := h.storageService.GetStorage(r.Context(), id)
	if err != nil || st == nil {
		WriteError(w, http.StatusNotFound, "storage record not found")
		return
	}
	WriteJSON(w, http.StatusOK, st)
}

func (h *StorageHandler) DeleteStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}
	if err := h.storageService.DeleteStorage(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *StorageHandler) StartStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}
	st, err := h.storageService.StartStorage(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, st)
}

func (h *StorageHandler) StopStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}
	if err := h.storageService.StopStorage(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
