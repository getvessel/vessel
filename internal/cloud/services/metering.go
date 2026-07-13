package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meterevent"
	"vessl.dev/vessl/internal/cloud/notifications"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/models"
)

type MeteringService interface {
	RecordUsage(teamID uint, deployments int, containerHours int, bandwidthGB int) error
}

type DefaultMeteringService struct {
	repo   repos.CloudRepo
	mailer *notifications.MailerService
}

func NewMeteringService(repo repos.CloudRepo, mailer *notifications.MailerService) *DefaultMeteringService {
	return &DefaultMeteringService{repo: repo, mailer: mailer}
}

func (s *DefaultMeteringService) RecordUsage(teamID uint, deployments int, containerHours int, bandwidthGB int) error {
	err := s.repo.LogUsage(&models.CloudUsageLog{
		WorkspaceID:    teamID,
		Deployments:    deployments,
		ContainerHours: containerHours,
		BandwidthGB:    bandwidthGB,
		ReportedAt:     time.Now(),
	})
	if err != nil {
		log.Printf("Failed to record usage in DB: %v", err)
	}

	team, err := s.repo.GetTeamByID(teamID)
	if err != nil || team == nil {
		return err
	}

	if team.StripeCustomerID == "" {
		return nil
	}

	if deployments > 0 {
		if err := s.reportToStripe(team.StripeCustomerID, "deployments_meter", deployments); err != nil {
			log.Printf("Failed to report deployment usage to Stripe for %d: %v", teamID, err)
		}
	}

	if containerHours > 0 {
		if err := s.reportToStripe(team.StripeCustomerID, "container_hours_meter", containerHours); err != nil {
			log.Printf("Failed to report container usage to Stripe for %d: %v", teamID, err)
		}
	}

	if bandwidthGB > 0 {
		if err := s.reportToStripe(team.StripeCustomerID, "bandwidth_gb_meter", bandwidthGB); err != nil {
			log.Printf("Failed to report bandwidth usage to Stripe for %d: %v", teamID, err)
		}
	}

	totalHours, _, err := s.repo.GetCurrentMonthUsage(teamID)
	if err == nil {
		ctx := context.Background()
		teamStrID := fmt.Sprintf("team_%d", teamID)

		hoursLimit := GetFeatures().GetMaxServers(teamStrID, team.Plan) * 730
		if hoursLimit > 0 {
			percentage := (totalHours * 100) / hoursLimit
			if percentage >= 80 && percentage <= 100 {
				owners, _ := s.repo.GetTeamOwners(teamID)
				for _, owner := range owners {
					_ = s.mailer.SendUsageAlertEmail(ctx, owner.UserEmail, team.Name, "Container Hours", totalHours, hoursLimit, percentage)
				}
			}
		}
	}

	return nil
}

func (s *DefaultMeteringService) reportToStripe(customerID string, eventName string, value int) error {
	params := &stripe.BillingMeterEventParams{
		EventName: stripe.String(eventName),
		Payload: map[string]string{
			"stripe_customer_id": customerID,
			"value":              fmt.Sprintf("%d", value),
		},
		Timestamp: stripe.Int64(time.Now().Unix()),
	}

	_, err := meterevent.New(params)
	return err
}
