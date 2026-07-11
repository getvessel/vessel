package services

import (
	"context"
	"log"
	"time"
)

// AuditEvent represents an action taken by a user/system
type AuditEvent struct {
	ID        string
	TeamID    string
	UserID    string
	Action    string
	Resource  string
	IPAddress string
	Timestamp time.Time
}

type AuditService struct {
	// db *repos.CloudDB
}

func NewAuditService() *AuditService {
	return &AuditService{}
}

// LogEvent securely records an audit event
func (s *AuditService) LogEvent(ctx context.Context, event AuditEvent) error {
	event.Timestamp = time.Now()

	// Mock: in production this inserts into an append-only Postgres table or separate data store
	log.Printf("[AUDIT] Team: %s | User: %s | Action: %s | Resource: %s",
		event.TeamID, event.UserID, event.Action, event.Resource)

	return nil
}
