package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/services"
)

type CanvasHandler struct {
	canvasService *services.CanvasService
}

func NewCanvasHandler(s *services.CanvasService) *CanvasHandler {
	return &CanvasHandler{canvasService: s}
}

func (h *CanvasHandler) ListCanvasSummaries(c echo.Context) error {
	summaries, err := h.canvasService.ListSummaries(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, summaries)
}

func (h *CanvasHandler) GetCanvasSummary(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	summary, err := h.canvasService.GetSummary(r.Context(), id)
	if err != nil || summary == nil {
		WriteError(w, http.StatusNotFound, "canvas summary not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, summary)
}

func (h *CanvasHandler) GetEnvironmentCanvas(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	canvas, err := h.canvasService.GetEnvironmentCanvas(r.Context(), id)
	if err != nil || canvas == nil {
		WriteError(w, http.StatusNotFound, "environment canvas not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, canvas)
}
