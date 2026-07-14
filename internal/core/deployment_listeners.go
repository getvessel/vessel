package core

import (
	"fmt"
	"log/slog"

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
	slog.Info("deployment completed", "serviceID", e.ServiceID, "status", e.Status)
}

func (l *DeploymentListeners) TriggerWebhook(e DeploymentCompleted) {
	slog.Info("triggering webhook", "projectID", e.ProjectID)
}

func (l *DeploymentListeners) Register() {
	On("deployment.completed", l.SendNotification)
	On("deployment.completed", l.UpdateAuditLog)
	On("deployment.completed", l.TriggerWebhook)
}
