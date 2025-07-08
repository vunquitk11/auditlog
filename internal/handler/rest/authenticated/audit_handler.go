package authenticated

import (
	"net/http"
	"strconv"

	"github.com/audit-log-service/internal/model"
	"github.com/gin-gonic/gin"
)

// GetAuditLogs handles GET /audit-logs and returns a list of audit logs based on filter criteria.
func (h *Handler) GetAuditLogs(c *gin.Context) {
	// Parse filter from query params
	criteria := model.FilterCriteria{
		ServiceName:  c.Query("service_name"),
		Username:     c.Query("username"),
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		ResourceID:   c.Query("resource_id"),
	}
	if page, ok := c.GetQuery("page"); ok {
		criteria.Page = atoiOrDefault(page, 1)
	}
	if pageSize, ok := c.GetQuery("page_size"); ok {
		criteria.PageSize = atoiOrDefault(pageSize, 20)
	}

	logs, err := h.auditLogCtrl.GetAuditLogs(c.Request.Context(), criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// atoiOrDefault converts a string to int, returning a default value if conversion fails.
func atoiOrDefault(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}
