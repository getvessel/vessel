package handlers

import (
	"github.com/labstack/echo/v4"
)

type BillingHandler struct{}

func NewBillingHandler() *BillingHandler {
	return &BillingHandler{}
}

func (h *BillingHandler) SetupCheckout(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "checkout_initialized"})
}
