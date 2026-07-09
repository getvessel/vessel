package team

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
