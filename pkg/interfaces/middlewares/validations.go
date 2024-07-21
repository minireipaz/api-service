package middlewares

import (
	"minireipaz/pkg/domain/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateWorkflow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var workflow models.Workflow
		if err := c.ShouldBindJSON(&workflow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			c.Abort()
			return
		}

		if strings.TrimSpace(workflow.WorkflowName) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow name is required"})
			c.Abort()
			return
		}

		if len(workflow.WorkflowName) < 3 || len(workflow.WorkflowName) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow name must be between 3 and 50 characters"})
			c.Abort()
			return
		}

		if strings.TrimSpace(workflow.DirectoryToSave) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Directory to save is required"})
			c.Abort()
			return
		}

		c.Set("workflow", workflow)
		c.Next()
	}
}
