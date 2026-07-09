package domain

import (
	"context"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service     *Service
	proxyReload func(ctx context.Context) error
}

type ProxyReloader interface {
	Reload(ctx context.Context) error
}

func NewHandler(service *Service, proxy ProxyReloader) *Handler {
	reload := func(ctx context.Context) error { return nil }
	if proxy != nil {
		reload = proxy.Reload
	}
	return &Handler{service: service, proxyReload: reload}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	domains, err := h.service.ListByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, domains)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	var d Config
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid domain payload")
		return
	}
	d.ProjectID = projectID
	if d.DomainName == "" {
		writeError(w, http.StatusBadRequest, "domain_name is required")
		return
	}

	if err := h.service.Create(r.Context(), &d); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.proxyReload(r.Context())
	writeJSON(w, http.StatusCreated, d)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.proxyReload(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
