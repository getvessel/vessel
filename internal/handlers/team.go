package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(s *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: s}
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	teams, err := h.teamService.ListTeamsByUser(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, teams)
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		WriteError(w, http.StatusBadRequest, "valid team name required")
		return
	}
	team, err := h.teamService.CreateTeam(r.Context(), req.Name, userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, team)
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return
	}
	team, err := h.teamService.GetTeam(r.Context(), id)
	if err != nil || team == nil {
		WriteError(w, http.StatusNotFound, "team not found")
		return
	}
	WriteJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return
	}
	if err := h.teamService.DeleteTeam(r.Context(), id, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return
	}
	members, err := h.teamService.ListMembers(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, members)
}

func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return
	}
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		WriteError(w, http.StatusBadRequest, "valid email required")
		return
	}
	inv, err := h.teamService.InviteMember(r.Context(), id, req.Email, req.Role)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, inv)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	targetUserID := r.PathValue("userId")
	if id == "" || targetUserID == "" {
		WriteError(w, http.StatusBadRequest, "missing team id or userId")
		return
	}
	if err := h.teamService.RemoveMember(r.Context(), id, targetUserID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) GetInvite(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		WriteError(w, http.StatusBadRequest, "missing invite token")
		return
	}
	inv, err := h.teamService.GetInvite(r.Context(), token)
	if err != nil || inv == nil {
		WriteError(w, http.StatusNotFound, "invite not found or expired")
		return
	}
	WriteJSON(w, http.StatusOK, inv)
}

func (h *TeamHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	token := r.PathValue("token")
	if token == "" {
		WriteError(w, http.StatusBadRequest, "missing invite token")
		return
	}
	if err := h.teamService.AcceptInvite(r.Context(), token, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "accepted"})
}
