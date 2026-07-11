package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
)

type BillingHandler struct {
	stripeWebhookSecret string
	paddleWebhookSecret string
}

func NewBillingHandler() *BillingHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &BillingHandler{
		stripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		paddleWebhookSecret: os.Getenv("PADDLE_WEBHOOK_SECRET"),
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

// HandlePaddleWebhook handles Paddle subscription webhooks
// @Summary Paddle Webhook
// @Description Receives billing events from Paddle
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/paddle/webhook [post]
func (h *BillingHandler) HandlePaddleWebhook(c echo.Context) error {
	const MaxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, MaxBodyBytes)

	if h.paddleWebhookSecret != "" {
		verifier := paddle.NewWebhookVerifier(h.paddleWebhookSecret)
		ok, err := verifier.Verify(c.Request())
		if err != nil {
			log.Printf("Paddle verification error: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Verification failed"})
		}
		if !ok {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid signature"})
		}
	}

	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Error reading request body"})
	}

	// Payload is verified, parse the JSON
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	eventType, _ := event["event_type"].(string)
	log.Printf("Received Paddle webhook event: %s", eventType)

	switch eventType {
	case "subscription.created", "subscription.updated":
		log.Println("Handling Paddle subscription update...")
		// TODO: Update cloud_subscriptions table

	case "subscription.canceled":
		log.Println("Handling Paddle subscription canceled...")
		// TODO: Mark as canceled in cloud_subscriptions table
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
}

type CheckoutRequest struct {
	PlanID    string `json:"plan_id"`
	ReturnURL string `json:"return_url"`
}

// CreateStripeCheckout creates a new Stripe Checkout session
// @Summary Create Stripe Checkout
// @Description Generates a checkout URL for subscriptions
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/stripe/checkout [post]
func (h *BillingHandler) CreateStripeCheckout(c echo.Context) error {
	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Mock mapping of plan names to Stripe Price IDs
	// In production, these should come from environment variables or a DB
	priceID := ""
	switch req.PlanID {
	case "hobby":
		priceID = os.Getenv("STRIPE_PRICE_HOBBY")
	case "pro":
		priceID = os.Getenv("STRIPE_PRICE_PRO")
	case "team":
		priceID = os.Getenv("STRIPE_PRICE_TEAM")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid plan ID"})
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(req.ReturnURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(req.ReturnURL + "?canceled=true"),
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("Stripe checkout error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create checkout session"})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}

// CreatePaddleCheckout creates a new Paddle Checkout session
// @Summary Create Paddle Checkout
// @Description Generates a checkout URL/Transaction for subscriptions
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/paddle/checkout [post]
func (h *BillingHandler) CreatePaddleCheckout(c echo.Context) error {
	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Currently returning a mock structure as Paddle v5 typically builds
	// transactions via API or relies on Paddle.js for checkout overlay.
	priceID := ""
	switch req.PlanID {
	case "hobby":
		priceID = os.Getenv("PADDLE_PRICE_HOBBY")
	case "pro":
		priceID = os.Getenv("PADDLE_PRICE_PRO")
	case "team":
		priceID = os.Getenv("PADDLE_PRICE_TEAM")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid plan ID"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"price_id": priceID,
		"status":   "ready_for_frontend_paddle_js",
	})
}
