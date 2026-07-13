package core

import (
	"fmt"
	"log"

	"vessl.dev/vessl/internal/models"
)

type DeploymentListeners struct {
	dispatcher *DispatcherService
}

func NewDeploymentListeners(dispatcher *DispatcherService) *DeploymentListeners {
	return &DeploymentListeners{dispatcher: dispatcher}
}

func (l *DeploymentListeners) SendNotification(e DeploymentCompleted) {
	commit := e.CommitHash
	if len(commit) > 7 {
		commit = commit[:7]
	}
	msg := fmt.Sprintf("Deploy %s: %s (%s)", e.Status, e.ServiceID, commit)
	notifEvent := &models.NotificationEvent{
		ProjectID: e.ProjectID,
		Level:     e.Status,
		Title:     "Deployment " + e.Status,
		Message:   msg,
		URL:       fmt.Sprintf("https://vessl.local/projects/%s/services/%s", e.ProjectID, e.ServiceID),
	}
	l.dispatcher.Dispatch(notifEvent)
}

func (l *DeploymentListeners) UpdateAuditLog(e DeploymentCompleted) {
	log.Printf("[Audit] Action: deployment.completed, ResourceID: %s, Status: %s", e.ServiceID, e.Status)
}

func (l *DeploymentListeners) TriggerWebhook(e DeploymentCompleted) {
	log.Printf("[Webhook] Triggering webhook for ProjectID: %s", e.ProjectID)
}

func (l *DeploymentListeners) Register() {
	On("deployment.completed", l.SendNotification)
	On("deployment.completed", l.UpdateAuditLog)
	On("deployment.completed", l.TriggerWebhook)
}
