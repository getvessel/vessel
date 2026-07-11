package services

import (
	"fmt"
	"log"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meterevent"
	"vessel.dev/vessel/internal/cloud/repos"
	"vessel.dev/vessel/internal/cloud/models"
)

type MeteringService interface {
	RecordUsage(teamID uint, deployments int, containerHours int, bandwidthGB int) error
	ReportToStripe(customerID string, eventName string, value int) error
}

type DefaultMeteringService struct{
	repo repos.CloudRepo
}

func NewMeteringService(repo repos.CloudRepo) *DefaultMeteringService {
	return &DefaultMeteringService{
		repo: repo,
	}
}

// RecordUsage stores usage internally in our database and conditionally pushes it to external billing providers
func (s *DefaultMeteringService) RecordUsage(teamID uint, deployments int, containerHours int, bandwidthGB int) error {
	// Save the raw metrics to Postgres table `cloud_usage_logs`
	err := s.repo.LogUsage(&models.CloudUsageLog{
		TeamID:         teamID,
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
		return err // Team not found or DB error
	}

	customerStripeID := team.StripeCustomerID
	if customerStripeID == "" {
		return nil // Not connected to Stripe
	}

	// Report metric events to Stripe if they are on a metered plan
	if deployments > 0 {
		err := s.ReportToStripe(customerStripeID, "deployments_meter", deployments)
		if err != nil {
			log.Printf("Failed to report deployment usage to Stripe for %d: %v", teamID, err)
		}
	}
	
	if containerHours > 0 {
		err := s.ReportToStripe(customerStripeID, "container_hours_meter", containerHours)
		if err != nil {
			log.Printf("Failed to report container usage to Stripe for %d: %v", teamID, err)
		}
	}

	if bandwidthGB > 0 {
		err := s.ReportToStripe(customerStripeID, "bandwidth_gb_meter", bandwidthGB)
		if err != nil {
			log.Printf("Failed to report bandwidth usage to Stripe for %d: %v", teamID, err)
		}
	}

	return nil
}

// ReportToStripe pushes a single usage record to Stripe's v2 metered billing API
func (s *DefaultMeteringService) ReportToStripe(customerID string, eventName string, value int) error {
	// Create a new metering event in Stripe
	params := &stripe.BillingMeterEventParams{
		EventName: stripe.String(eventName),
		Payload: map[string]string{
			"stripe_customer_id": customerID,
			"value":              fmt.Sprintf("%d", value),
		},
		Timestamp: stripe.Int64(time.Now().Unix()),
	}

	_, err := meterevent.New(params)
	if err != nil {
		return err
	}

	return nil
}
