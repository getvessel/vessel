package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Migrate(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar_url TEXT NOT NULL DEFAULT '',
			preferred_region TEXT NOT NULL DEFAULT 'local',
			owner_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			user_id TEXT NOT NULL,
			user_email TEXT NOT NULL DEFAULT '',
			role TEXT NOT NULL DEFAULT 'Member',
			joined_at TEXT NOT NULL,
			UNIQUE(team_id, user_id)
		);
		CREATE TABLE IF NOT EXISTS team_invites (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			email TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'Member',
			token TEXT NOT NULL UNIQUE,
			invited_by TEXT NOT NULL,
			expires_at TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
	`)
	return err
}

func (r *SQLiteRepository) CreateTeam(ctx context.Context, team *Team) error {
	if team.ID == "" {
		team.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if team.CreatedAt.IsZero() {
		team.CreatedAt = now
	}
	team.UpdatedAt = team.CreatedAt

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `INSERT INTO teams (id, name, avatar_url, preferred_region, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		team.ID, team.Name, team.AvatarURL, team.PreferredRegion, team.OwnerID, team.CreatedAt.Format(time.RFC3339), team.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}

	ownerMember := &TeamMember{
		ID:        uuid.NewString(),
		TeamID:    team.ID,
		UserID:    team.OwnerID,
		UserEmail: "",
		Role:      "Owner",
		JoinedAt:  now,
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO team_members (id, team_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		ownerMember.ID, ownerMember.TeamID, ownerMember.UserID, ownerMember.UserEmail, ownerMember.Role, ownerMember.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert team owner: %w", err)
	}

	return tx.Commit()
}

func (r *SQLiteRepository) GetTeamByID(ctx context.Context, id string) (*Team, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var t Team
	var createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM teams WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.AvatarURL, &t.PreferredRegion, &t.OwnerID, &createdStr, &updatedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get team %s: %w", id, err)
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &t, nil
}

func (r *SQLiteRepository) ListTeamsByUser(ctx context.Context, userID string) ([]*Team, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `SELECT t.id, t.name, t.avatar_url, t.preferred_region, t.owner_id, t.created_at, t.updated_at
	          FROM teams t
	          JOIN team_members m ON t.id = m.team_id
	          WHERE m.user_id = ? ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list teams by user: %w", err)
	}
	defer rows.Close()

	var list []*Team
	for rows.Next() {
		var t Team
		var createdStr, updatedStr string
		if err := rows.Scan(&t.ID, &t.Name, &t.AvatarURL, &t.PreferredRegion, &t.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &t)
	}
	return list, nil
}

func (r *SQLiteRepository) UpdateTeam(ctx context.Context, team *Team) error {
	team.UpdatedAt = time.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `UPDATE teams SET name = ?, avatar_url = ?, preferred_region = ?, updated_at = ? WHERE id = ?`,
		team.Name, team.AvatarURL, team.PreferredRegion, team.UpdatedAt.Format(time.RFC3339), team.ID)
	return err
}

func (r *SQLiteRepository) DeleteTeam(ctx context.Context, id, ownerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.ExecContext(ctx, `DELETE FROM teams WHERE id = ? AND owner_id = ?`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete team: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("team not found or unauthorized (must be owner)")
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM team_members WHERE team_id = ?`, id)
	_, _ = r.db.ExecContext(ctx, `DELETE FROM team_invites WHERE team_id = ?`, id)
	return nil
}

func (r *SQLiteRepository) AddMember(ctx context.Context, member *TeamMember) error {
	if member.ID == "" {
		member.ID = uuid.NewString()
	}
	if member.JoinedAt.IsZero() {
		member.JoinedAt = time.Now().UTC()
	}
	if member.Role == "" {
		member.Role = "Member"
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `INSERT INTO team_members (id, team_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		member.ID, member.TeamID, member.UserID, member.UserEmail, member.Role, member.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("add team member: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) RemoveMember(ctx context.Context, teamID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.ExecContext(ctx, `DELETE FROM team_members WHERE team_id = ? AND user_id = ? AND role != 'Owner'`, teamID, userID)
	if err != nil {
		return fmt.Errorf("remove team member: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("member not found or cannot remove team Owner")
	}
	return nil
}

func (r *SQLiteRepository) GetMember(ctx context.Context, teamID, userID string) (*TeamMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var m TeamMember
	var joinedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, team_id, user_id, user_email, role, joined_at FROM team_members WHERE team_id = ? AND user_id = ?`, teamID, userID).
		Scan(&m.ID, &m.TeamID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get team member: %w", err)
	}
	m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
	return &m, nil
}

func (r *SQLiteRepository) ListMembers(ctx context.Context, teamID string) ([]*TeamMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT id, team_id, user_id, user_email, role, joined_at FROM team_members WHERE team_id = ? ORDER BY joined_at ASC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("list team members: %w", err)
	}
	defer rows.Close()

	var list []*TeamMember
	for rows.Next() {
		var m TeamMember
		var joinedStr string
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr); err != nil {
			return nil, err
		}
		m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
		list = append(list, &m)
	}
	return list, nil
}

func (r *SQLiteRepository) CreateInvite(ctx context.Context, invite *TeamInvite) error {
	if invite.ID == "" {
		invite.ID = uuid.NewString()
	}
	if invite.Token == "" {
		invite.Token = uuid.NewString()
	}
	now := time.Now().UTC()
	if invite.CreatedAt.IsZero() {
		invite.CreatedAt = now
	}
	if invite.ExpiresAt.IsZero() {
		invite.ExpiresAt = now.Add(7 * 24 * time.Hour)
	}
	if invite.Role == "" {
		invite.Role = "Member"
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `INSERT INTO team_invites (id, team_id, email, role, token, invited_by, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		invite.ID, invite.TeamID, invite.Email, invite.Role, invite.Token, invite.InvitedBy, invite.ExpiresAt.Format(time.RFC3339), invite.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create team invite: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetInviteByToken(ctx context.Context, token string) (*TeamInvite, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var inv TeamInvite
	var expStr, createdStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, team_id, email, role, token, invited_by, expires_at, created_at FROM team_invites WHERE token = ?`, token).
		Scan(&inv.ID, &inv.TeamID, &inv.Email, &inv.Role, &inv.Token, &inv.InvitedBy, &expStr, &createdStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get team invite: %w", err)
	}
	inv.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
	inv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	return &inv, nil
}

func (r *SQLiteRepository) DeleteInvite(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `DELETE FROM team_invites WHERE id = ?`, id)
	return err
}
