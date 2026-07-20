package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type ScheduledTaskHandler struct {
	scheduledTaskService *services.ScheduledTaskService
}

func NewScheduledTaskHandler(s *services.ScheduledTaskService) *ScheduledTaskHandler {
	return &ScheduledTaskHandler{scheduledTaskService: s}
}

// @Summary ListProjectScheduledTasks endpoint
// @Description ListProjectScheduledTasks endpoint
// @Tags ScheduledTasks
// @Accept json
// @Produce json
// @Router /scheduled-tasks [get]
func (h *ScheduledTaskHandler) ListProjectScheduledTasks(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	serviceID := c.QueryParam("serviceId")

	var tasks []models.ScheduledTask
	var err error

	if serviceID != "" {
		tasks, err = h.scheduledTaskService.ListScheduledTasksByService(c.Request().Context(), serviceID)
	} else {
		tasks, err = h.scheduledTaskService.ListScheduledTasksByProject(c.Request().Context(), projectID)
	}

	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", tasks)
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags ScheduledTasks
// @Accept json
// @Produce json
// @Param request body models.ScheduledTask true "Payload"
// @Router /scheduled-tasks [post]
func (h *ScheduledTaskHandler) Create(c echo.Context) error {
	var j models.ScheduledTask
	if err := c.Bind(&j); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	created, err := h.scheduledTaskService.CreateScheduledTask(c.Request().Context(), &j)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

// @Summary Get ScheduledTask
// @Description Get ScheduledTask
// @Tags ScheduledTasks
// @Accept json
// @Produce json
// @Param id path string true "ScheduledTask ID"
// @Router /scheduled-tasks/{id} [get]
func (h *ScheduledTaskHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	j, err := h.scheduledTaskService.GetScheduledTask(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	return utils.Success(c, "Operation successful", j)
}

// @Summary Delete ScheduledTask
// @Description Delete ScheduledTask
// @Tags ScheduledTasks
// @Accept json
// @Produce json
// @Param id path string true "ScheduledTask ID"
// @Router /scheduled-tasks/{id} [delete]
func (h *ScheduledTaskHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.scheduledTaskService.DeleteScheduledTask(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary Run endpoint
// @Description Run endpoint
// @Tags ScheduledTasks
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /scheduled-tasks/{id}/trigger [post]
func (h *ScheduledTaskHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	out, err := h.scheduledTaskService.ExecuteScheduledTask(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "executed", "output": out})
}
