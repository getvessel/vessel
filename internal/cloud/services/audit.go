package services

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

func (s *AuditService) LogEvent(teamID, userID, action string) error {
	return nil
}
