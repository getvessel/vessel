package models

import "time"

type Team struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	AvatarURL       string    `json:"avatarUrl,omitempty"`
	PreferredRegion string    `json:"preferredRegion,omitempty"`
	OwnerID         string    `json:"ownerId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type TeamMember struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	UserID    string    `json:"userId"`
	UserEmail string    `json:"userEmail"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
}

type TeamInvite struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	InvitedBy string    `json:"invitedBy"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateTeamRequest struct {
	Name            string `json:"name"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type InviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AcceptInviteRequest struct {
	Token string `json:"token"`
}

type GetTeamResponse struct {
	Team    *Team         `json:"team"`
	Members []*TeamMember `json:"members"`
}

type Workspace struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	AvatarURL       string    `json:"avatarUrl,omitempty"`
	PreferredRegion string    `json:"preferredRegion,omitempty"`
	OwnerID         string    `json:"ownerId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type TrustedDomain struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Domain    string    `json:"domain"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type SSHKey struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Name      string    `json:"name"`
	PublicKey string    `json:"publicKey"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuditLog struct {
	ID            string    `json:"id"`
	TeamID        string    `json:"teamId"`
	ProjectID     string    `json:"projectId,omitempty"`
	EnvironmentID string    `json:"environmentId,omitempty"`
	Action        string    `json:"action"`
	Actor         string    `json:"actor"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CreateWorkspaceRequest struct {
	Name            string `json:"name"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type UpdateWorkspaceRequest struct {
	Name            string `json:"name,omitempty"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type CreateTrustedDomainRequest struct {
	Domain string `json:"domain"`
	Role   string `json:"role"`
}

type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}
