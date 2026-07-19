package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type DNSHandler struct {
	dnsService *services.DNSService
}

func NewDNSHandler(dnsService *services.DNSService) *DNSHandler {
	return &DNSHandler{dnsService: dnsService}
}

// @Summary Create a DNS Record
// @Description Create a generic DNS record (A, CNAME, TXT, etc.)
// @Tags DNS
// @Accept json
// @Produce json
// @Param request body models.CreateDNSRecordRequest true "Payload"
// @Router /dns [post]
func (h *DNSHandler) Create(c echo.Context) error {
	var req models.CreateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.CreateRecord(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Created(c, "DNS record created successfully", record)
}

// @Summary List DNS Records
// @Description List all DNS records for a specific domain
// @Tags DNS
// @Accept json
// @Produce json
// @Param domain query string true "Domain Name"
// @Router /dns [get]
func (h *DNSHandler) List(c echo.Context) error {
	domain := c.QueryParam("domain")

	records, err := h.dnsService.ListByDomain(c.Request().Context(), domain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS records fetched successfully", records)
}

// @Summary Update a DNS Record
// @Description Update an existing DNS record
// @Tags DNS
// @Accept json
// @Produce json
// @Param id path string true "Record ID"
// @Param request body models.UpdateDNSRecordRequest true "Payload"
// @Router /dns/{id} [put]
func (h *DNSHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	var req models.UpdateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.UpdateRecord(c.Request().Context(), id, &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS record updated successfully", record)
}

// @Summary Delete a DNS Record
// @Description Delete an existing DNS record
// @Tags DNS
// @Accept json
// @Produce json
// @Param id path string true "Record ID"
// @Router /dns/{id} [delete]
func (h *DNSHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	if err := h.dnsService.DeleteRecord(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
