package service_var

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo Repository
	svc  ServiceRepository
}

func NewHandler(repo Repository, svc ServiceRepository) *Handler {
	return &Handler{repo: repo, svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	list, err := h.repo.ListByService(r.Context(), serviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	var req CreateServiceVarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	svc, err := h.svc.GetByID(r.Context(), serviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if svc == nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	v := &Variable{
		ServiceID:     serviceID,
		ProjectID:     svc.ProjectID,
		EnvironmentID: svc.EnvironmentID,
		Key:           req.Key,
		Value:         req.Value,
		IsSecret:      req.IsSecret,
	}

	if err := h.repo.Create(r.Context(), v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	id := r.PathValue("id")

	var req UpdateServiceVarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v := &Variable{
		ID:        id,
		ServiceID: serviceID,
		Key:       req.Key,
		Value:     req.Value,
		IsSecret:  req.IsSecret,
	}

	if err := h.repo.Update(r.Context(), v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	id := r.PathValue("id")

	if err := h.repo.Delete(r.Context(), id, serviceID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
