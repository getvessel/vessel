package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type BillingHandler struct {
	stripeWebhookSecret string
}

func NewBillingHandler() *BillingHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &BillingHandler{
		stripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}

// HandleStripeWebhook handles Stripe subscription webhooks
// @Summary Stripe Webhook
// @Description Receives billing events from Stripe to update subscription status
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/stripe/webhook [post]
func (h *BillingHandler) HandleStripeWebhook(c echo.Context) error {
	const MaxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Error reading request body"})
	}

	sigHeader := c.Request().Header.Get("Stripe-Signature")

	var event stripe.Event

	// Only verify signature if secret is configured (useful for local dev)
	if h.stripeWebhookSecret != "" {
		event, err = webhook.ConstructEvent(payload, sigHeader, h.stripeWebhookSecret)
		if err != nil {
			log.Printf("Stripe signature verification failed: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid signature"})
		}
	} else {
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("Failed to parse webhook body json: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
		}
	}

	log.Printf("Received Stripe webhook event: %s", event.Type)

	switch event.Type {
	case "customer.subscription.created", "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid subscription payload"})
		}

		status := subscription.Status
		customerID := subscription.Customer.ID
		planID := subscription.Items.Data[0].Price.ID

		log.Printf("Subscription Update | Customer: %s | Status: %s | Plan: %s", customerID, status, planID)
		// TODO: Update cloud_subscriptions table via CloudDB

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid subscription payload"})
		}

		log.Printf("Subscription Deleted | Customer: %s", subscription.Customer.ID)
		// TODO: Mark as canceled in cloud_subscriptions table
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
}

// HandlePaystackWebhook handles Paystack subscription webhooks
// @Summary Paystack Webhook
// @Description Receives billing events from Paystack
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/paystack/webhook [post]
func (h *BillingHandler) HandlePaystackWebhook(c echo.Context) error {
	// Paystack signature verification requires HMAC-SHA512
	// For simplicity, we are leaving the stub but it would be similarly structured
	log.Println("Received Paystack webhook event")
	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
}
