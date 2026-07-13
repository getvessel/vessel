package handlers

import (
	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/cloud/utils"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// @Summary Get User Profile
// @Description Fetch current user details
// @Tags Cloud-Users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /cloud/users/me [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	return utils.Success(c, "Success", map[string]interface{}{
		"id":    "usr_123",
		"email": "user@example.com",
		"teams": []map[string]string{
			{"id": "team_1", "name": "Personal Team", "role": "owner"},
		},
	})
}
