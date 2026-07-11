package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BillingHandler struct {
	// db *repos.CloudDB
}

func NewBillingHandler() *BillingHandler {
	return &BillingHandler{}
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
	// In production, you would verify the Stripe signature:
	// sigHeader := c.Request().Header.Get("Stripe-Signature")

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read body"})
	}

	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	eventType, _ := event["type"].(string)
	log.Printf("Received Stripe webhook event: %s", eventType)

	switch eventType {
	case "customer.subscription.created", "customer.subscription.updated":
		// TODO: Parse subscription status and update `cloud_subscriptions` table
		log.Println("Handling subscription update...")
	case "customer.subscription.deleted":
		// TODO: Mark subscription as canceled in `cloud_subscriptions` table
		log.Println("Handling subscription deletion...")
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
	// Paystack signature verification would happen here
	log.Println("Received Paystack webhook event")
	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
}
