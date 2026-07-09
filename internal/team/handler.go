package team

import (
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/user"
)

type UserProvider interface {
	GetUserByEmail(email string) (*user.User, error)
}

type ClaimsExtractor func(r *http.Request) (userID, email, role string)

type Handler struct {
	repo          Repository
	userProvider  UserProvider
	extractClaims ClaimsExtractor
}

func NewHandler(repo Repository, userProvider UserProvider, extractClaims ClaimsExtractor) *Handler {
	return &Handler{
		repo:          repo,
		userProvider:  userProvider,
		extractClaims: extractClaims,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teams, err := h.repo.ListTeamsByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, teams)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "team name is required")
		return
	}

	team := &Team{
		Name:            req.Name,
		AvatarURL:       req.AvatarURL,
		PreferredRegion: req.PreferredRegion,
		OwnerID:         userID,
	}
	if err := h.repo.CreateTeam(r.Context(), team); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, team)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _, role := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	team, err := h.repo.GetTeamByID(r.Context(), teamID)
	if err != nil || team == nil {
		writeError(w, http.StatusNotFound, "team not found")
		return
	}

	member, _ := h.repo.GetMember(r.Context(), teamID, userID)
	if member == nil && role != "admin" {
		writeError(w, http.StatusForbidden, "not a member of this team")
		return
	}

	members, _ := h.repo.ListMembers(r.Context(), teamID)
	writeJSON(w, http.StatusOK, GetTeamResponse{
		Team:    team,
		Members: members,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	if err := h.repo.DeleteTeam(r.Context(), teamID, userID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("id")
	members, err := h.repo.ListMembers(r.Context(), teamID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *Handler) InviteMember(w http.ResponseWriter, r *http.Request) {
	userID, email, role := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	callerMember, _ := h.repo.GetMember(r.Context(), teamID, userID)
	if callerMember == nil || (callerMember.Role != "Owner" && callerMember.Role != "Admin") {
		if role != "admin" {
			writeError(w, http.StatusForbidden, "only Owner or Admin can invite members")
			return
		}
	}

	var req InviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid invite payload")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Role == "" {
		req.Role = "Member"
	}

	existingUser, _ := h.userProvider.GetUserByEmail(req.Email)
	if existingUser != nil {
		member := &TeamMember{
			TeamID:    teamID,
			UserID:    existingUser.ID,
			UserEmail: existingUser.Email,
			Role:      req.Role,
		}
		if err := h.repo.AddMember(r.Context(), member); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"status": "added_to_team",
			"member": member,
		})
		return
	}

	invite := &TeamInvite{
		TeamID:    teamID,
		Email:     req.Email,
		Role:      req.Role,
		InvitedBy: email,
	}
	if err := h.repo.CreateInvite(r.Context(), invite); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"status": "invitation_sent",
		"invite": invite,
	})
}

func (h *Handler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, _, role := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	targetUserID := r.PathValue("userId")

	callerMember, _ := h.repo.GetMember(r.Context(), teamID, userID)
	if callerMember == nil || (callerMember.Role != "Owner" && callerMember.Role != "Admin" && userID != targetUserID) {
		if role != "admin" {
			writeError(w, http.StatusForbidden, "unauthorized to remove this member")
			return
		}
	}

	if err := h.repo.RemoveMember(r.Context(), teamID, targetUserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetInvite(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	inv, err := h.repo.GetInviteByToken(r.Context(), token)
	if err != nil || inv == nil {
		writeError(w, http.StatusNotFound, "invitation not found or expired")
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

func (h *Handler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userID, email, _ := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	token := r.PathValue("token")
	inv, err := h.repo.GetInviteByToken(r.Context(), token)
	if err != nil || inv == nil {
		writeError(w, http.StatusNotFound, "invitation not found or expired")
		return
	}

	member := &TeamMember{
		TeamID:    inv.TeamID,
		UserID:    userID,
		UserEmail: email,
		Role:      inv.Role,
	}
	if err := h.repo.AddMember(r.Context(), member); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = h.repo.DeleteInvite(r.Context(), inv.ID)

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "accepted",
		"member": member,
	})
}
