package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

// @Summary ListVolumes endpoint
// @Description ListVolumes endpoint
// @Tags ServiceVolumes
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /apps/{id}/volumes [get]
func (h *AppHandler) ListVolumes(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "Service not found")
	}

	list, err := h.appService.ListVolumes(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to list volumes")
	}
	return utils.Success(c, "Operation successful", list)
}

// @Summary CreateVolume endpoint
// @Description CreateVolume endpoint
// @Tags ServiceVolumes
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body models.ServiceVolume true "Payload"
// @Router /apps/{id}/volumes [post]
func (h *AppHandler) CreateVolume(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "Service not found")
	}

	var req models.ServiceVolume
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	req.ServiceID = serviceID
	created, err := h.appService.CreateVolume(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create volume")
	}
	return utils.Success(c, "Operation successful", created)
}

// @Summary DeleteVolume endpoint
// @Description DeleteVolume endpoint
// @Tags ServiceVolumes
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param volumeId path string true "volumeId"
// @Router /apps/{id}/volumes/{volumeId} [delete]
func (h *AppHandler) DeleteVolume(c echo.Context) error {
	serviceID := c.Param("id")
	volumeID := c.Param("volumeId")
	if serviceID == "" || volumeID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId and volumeId are required")
	}

	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "Service not found")
	}

	err = h.appService.DeleteVolume(c.Request().Context(), volumeID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to delete volume")
	}
	return utils.Success(c, "Operation successful", nil)
}
