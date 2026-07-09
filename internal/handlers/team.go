package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *TeamHandler) List(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	teams, err := h.teamService.ListTeamsByUser(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, teams)
}

func (h *TeamHandler) Create(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		WriteError(w, http.StatusBadRequest, "valid team name required")
		return nil
	}
	team, err := h.teamService.CreateTeam(r.Context(), req.Name, userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, team)
}

func (h *TeamHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return nil
	}
	team, err := h.teamService.GetTeam(r.Context(), id)
	if err != nil || team == nil {
		WriteError(w, http.StatusNotFound, "team not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) Delete(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return nil
	}
	if err := h.teamService.DeleteTeam(r.Context(), id, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) ListMembers(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return nil
	}
	members, err := h.teamService.ListMembers(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, members)
}

func (h *TeamHandler) InviteMember(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing team id")
		return nil
	}
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		WriteError(w, http.StatusBadRequest, "valid email required")
		return nil
	}
	inv, err := h.teamService.InviteMember(r.Context(), id, req.Email, req.Role)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, inv)
}

func (h *TeamHandler) RemoveMember(c echo.Context) error {
	id := c.Param("id")
	targetUserID := c.Param("userId")
	if id == "" || targetUserID == "" {
		WriteError(w, http.StatusBadRequest, "missing team id or userId")
		return nil
	}
	if err := h.teamService.RemoveMember(r.Context(), id, targetUserID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) GetInvite(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		WriteError(w, http.StatusBadRequest, "missing invite token")
		return nil
	}
	inv, err := h.teamService.GetInvite(r.Context(), token)
	if err != nil || inv == nil {
		WriteError(w, http.StatusNotFound, "invite not found or expired")
		return nil
	}
	WriteJSON(w, http.StatusOK, inv)
}

func (h *TeamHandler) AcceptInvite(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	token := c.Param("token")
	if token == "" {
		WriteError(w, http.StatusBadRequest, "missing invite token")
		return nil
	}
	if err := h.teamService.AcceptInvite(r.Context(), token, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "accepted"})
}
