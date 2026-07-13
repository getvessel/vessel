package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"
)

type CreateTrustedDomainRequest struct {
	Domain string `json:"domain"`
}

type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

// @Summary ListTrustedDomains endpoint
// @Description ListTrustedDomains endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/trusted-domains [get]
func (h *WorkspaceHandler) ListTrustedDomains(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	domains, err := h.workspaceService.ListTrustedDomains(c.Request().Context(), workspaceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", domains)
}

// @Summary CreateTrustedDomain endpoint
// @Description CreateTrustedDomain endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Param request body handlers.CreateTrustedDomainRequest true "Payload"
// @Router /workspaces/{workspaceId}/trusted-domains [post]
func (h *WorkspaceHandler) CreateTrustedDomain(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	var payload CreateTrustedDomainRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	td, err := h.workspaceService.AddTrustedDomain(c.Request().Context(), workspaceID, payload.Domain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", td)
}

// @Summary DeleteTrustedDomain endpoint
// @Description DeleteTrustedDomain endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
func (h *WorkspaceHandler) DeleteTrustedDomain(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.workspaceService.DeleteTrustedDomain(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListSSHKeys endpoint
// @Description ListSSHKeys endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/ssh-keys [get]
func (h *WorkspaceHandler) ListSSHKeys(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	keys, err := h.workspaceService.ListSSHKeys(c.Request().Context(), workspaceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", keys)
}

// @Summary CreateSSHKey endpoint
// @Description CreateSSHKey endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Param request body handlers.CreateSSHKeyRequest true "Payload"
// @Router /workspaces/{workspaceId}/ssh-keys [post]
func (h *WorkspaceHandler) CreateSSHKey(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	var payload CreateSSHKeyRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	key, err := h.workspaceService.AddSSHKey(c.Request().Context(), workspaceID, payload.Name, payload.PublicKey)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", key)
}

// @Summary DeleteSSHKey endpoint
// @Description DeleteSSHKey endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
func (h *WorkspaceHandler) DeleteSSHKey(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.workspaceService.DeleteSSHKey(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListAuditLogs endpoint
// @Description ListAuditLogs endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/audit-logs [get]
func (h *WorkspaceHandler) ListAuditLogs(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	logs, total, err := h.workspaceService.ListAuditLogs(c.Request().Context(), workspaceID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, "Operation successful", logs, total, page, limit)
}
