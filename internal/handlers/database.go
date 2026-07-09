package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type DatabaseHandler struct {
	databaseService *services.DatabaseService
}

func NewDatabaseHandler(s *services.DatabaseService) *DatabaseHandler {
	return &DatabaseHandler{databaseService: s}
}

func (h *DatabaseHandler) ListDatabases(w http.ResponseWriter, r *http.Request) {
	databases, err := h.databaseService.ListDatabases(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if databases == nil {
		databases = []*models.Database{}
	}
	WriteJSON(w, http.StatusOK, databases)
}

func (h *DatabaseHandler) CreateDatabase(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid database configuration payload")
		return
	}
	db, err := h.databaseService.CreateDatabaseFromRequest(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, db)
}

func (h *DatabaseHandler) GetDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}
	db, err := h.databaseService.GetDatabase(r.Context(), id)
	if err != nil || db == nil {
		WriteError(w, http.StatusNotFound, "database not found")
		return
	}
	WriteJSON(w, http.StatusOK, db)
}

func (h *DatabaseHandler) DeleteDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}
	if err := h.databaseService.DeleteDatabase(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *DatabaseHandler) StartDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}
	db, err := h.databaseService.StartDatabase(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, db)
}

func (h *DatabaseHandler) StopDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}
	if err := h.databaseService.StopDatabase(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
