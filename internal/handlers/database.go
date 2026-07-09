package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *DatabaseHandler) ListDatabases(c echo.Context) error {
	databases, err := h.databaseService.ListDatabases(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	if databases == nil {
		databases = []*models.Database{}
	}
	WriteJSON(w, http.StatusOK, databases)
}

func (h *DatabaseHandler) CreateDatabase(c echo.Context) error {
	var req models.CreateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	db, err := h.databaseService.CreateDatabaseFromRequest(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, db)
}

func (h *DatabaseHandler) GetDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return nil
	}
	db, err := h.databaseService.GetDatabase(r.Context(), id)
	if err != nil || db == nil {
		WriteError(w, http.StatusNotFound, "database not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, db)
}

func (h *DatabaseHandler) DeleteDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return nil
	}
	if err := h.databaseService.DeleteDatabase(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *DatabaseHandler) StartDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return nil
	}
	db, err := h.databaseService.StartDatabase(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, db)
}

func (h *DatabaseHandler) StopDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing database id parameter")
		return nil
	}
	if err := h.databaseService.StopDatabase(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
