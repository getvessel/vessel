package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MeteringHandler struct {
}

func NewMeteringHandler() *MeteringHandler {
	return &MeteringHandler{}
}

type UsageReport struct {
	TeamID         string `json:"team_id"`
	Deployments    int    `json:"deployments"`
	ContainerHours int    `json:"container_hours"`
	BandwidthGB    int    `json:"bandwidth_gb"`
}

// ReportUsage handles incoming usage metrics from connected agents
// @Summary Report Usage Metrics
// @Description Receives telemetry from Vossel Daemons for billing
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/usage/report [post]
func (h *MeteringHandler) ReportUsage(c echo.Context) error {
	var req UsageReport
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid usage report"})
	}

	// TODO: Validate agent token sending this report
	// TODO: Store metrics in PostgreSQL `cloud_usage_logs` table
	// TODO: Push line items to Stripe/Paddle usage-based billing APIs if threshold met

	log.Printf("Received usage report for team %s: %d deploys, %d hours", req.TeamID, req.Deployments, req.ContainerHours)

	return c.JSON(http.StatusOK, map[string]string{"status": "recorded"})
}
