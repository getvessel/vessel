package canvas

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) ListCanvasSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.repo.ListCanvasSummaries(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summaries)
}

func (h *Handler) GetCanvasSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	summary, err := h.repo.GetCanvasSummary(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *Handler) GetEnvironmentCanvas(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	canvas, err := h.repo.GetEnvironmentCanvas(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, canvas)
}
