package reports

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	protheticUserRole = "prothetic_user"
)

func (h *Handler) GetReport(c *gin.Context) {
	roles, ok := c.Get("roles")
	if !ok {
		c.JSON(http.StatusForbidden, "role required")
		return
	}

	if !checkRole(roles, protheticUserRole) {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	report := h.Getter.GetReport()
	c.JSON(http.StatusOK, report)
}

func checkRole(rawRoles any, requiredRole string) bool {
	roles, ok := rawRoles.([]string)
	if !ok {
		return false
	}

	for _, role := range roles {
		if role == requiredRole {
			return true
		}
	}

	return false
}
