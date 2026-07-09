package handlers

import (
	"net/http"

	"vessel.dev/vessel/internal/services"
)

type CanvasHandler struct {
	canvasService *services.CanvasService
}

func NewCanvasHandler(s *services.CanvasService) *CanvasHandler {
	return &CanvasHandler{canvasService: s}
}

func (h *CanvasHandler) ListCanvasSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.canvasService.ListSummaries(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, summaries)
}

func (h *CanvasHandler) GetCanvasSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	summary, err := h.canvasService.GetSummary(r.Context(), id)
	if err != nil || summary == nil {
		WriteError(w, http.StatusNotFound, "canvas summary not found")
		return
	}
	WriteJSON(w, http.StatusOK, summary)
}

func (h *CanvasHandler) GetEnvironmentCanvas(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	canvas, err := h.canvasService.GetEnvironmentCanvas(r.Context(), id)
	if err != nil || canvas == nil {
		WriteError(w, http.StatusNotFound, "environment canvas not found")
		return
	}
	WriteJSON(w, http.StatusOK, canvas)
}
