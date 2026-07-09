package team

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
