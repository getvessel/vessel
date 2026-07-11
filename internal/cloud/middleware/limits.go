package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/repos"
	"vessel.dev/vessel/internal/cloud/services"
)

// DeploymentRateLimiter intercepts deployment requests to check if the team has exceeded their hourly limit
func DeploymentRateLimiter(repo repos.CloudRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// In production, this is set by earlier AuthMiddleware (e.g. JWT claims)
			// Mocking team ID logic until real auth is plugged in, but DB queries use it
			teamID := uint(1) 
			teamStrID := "team_1"
			plan := "hobby"
			
			team, err := repo.GetTeamByID(teamID)
			if err == nil && team != nil {
				plan = team.Plan
			}

			limit := services.GetFeatures().GetDeploymentRateLimit(teamStrID, plan)

			currentUsage, _ := repo.GetDeploymentsInLastHour(teamID)

			if int(currentUsage) >= limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Deployment rate limit exceeded for your tier",
				})
			}

			return next(c)
		}
	}
}

// SeatLimitGuard intercepts server connection (BYOS) requests to check if they can add another server
func SeatLimitGuard(repo repos.CloudRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			teamID := uint(1)
			teamStrID := "team_1"
			plan := "hobby"

			team, err := repo.GetTeamByID(teamID)
			if err == nil && team != nil {
				plan = team.Plan
			}

			limit := services.GetFeatures().GetMaxServers(teamStrID, plan)

			currentServers, _ := repo.GetActiveServerCount(teamID)

			if int(currentServers) >= limit {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Bring Your Own Server (BYOS) seat limit reached. Please upgrade your plan.",
				})
			}

			return next(c)
		}
	}
}
