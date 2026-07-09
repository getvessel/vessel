package team

import "context"

type Repository interface {
	CreateTeam(ctx context.Context, team *Team) error
	GetTeamByID(ctx context.Context, id string) (*Team, error)
	ListTeamsByUser(ctx context.Context, userID string) ([]*Team, error)
	UpdateTeam(ctx context.Context, team *Team) error
	DeleteTeam(ctx context.Context, id, ownerID string) error

	AddMember(ctx context.Context, member *TeamMember) error
	RemoveMember(ctx context.Context, teamID, userID string) error
	GetMember(ctx context.Context, teamID, userID string) (*TeamMember, error)
	ListMembers(ctx context.Context, teamID string) ([]*TeamMember, error)

	CreateInvite(ctx context.Context, invite *TeamInvite) error
	GetInviteByToken(ctx context.Context, token string) (*TeamInvite, error)
	DeleteInvite(ctx context.Context, id string) error
}
