package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/cloud/repositories"
)

func RequireTeamRole(repo repos.CloudRepo, allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("cloud_user").(*CloudClaims)
			if !ok || claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			teamIDStr := c.Param("id")
			if teamIDStr == "" {
				teamIDStr = c.QueryParam("team_id")
			}
			teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
			}

			member, err := repo.GetTeamMember(uint(teamID), claims.UserID())
			if err != nil || member == nil {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: not a member of this team"})
			}

			if len(allowedRoles) > 0 {
				roleAllowed := false
				for _, role := range allowedRoles {
					if member.Role == role {
						roleAllowed = true
						break
					}
				}
				if !roleAllowed {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: insufficient permissions"})
				}
			}

			c.Set("team_member", member)
			return next(c)
		}
	}
}
