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
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON))
			c.Abort()
			return
		}

		if strings.TrimSpace(workflow.Sub) == "" {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserSubMustBe))
			c.Abort()
			return
		}
		if len(workflow.Sub) > 50 { // || !regexp.MustCompile(`^\d+$`).MatchString(workflow.Sub)
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserSubMustBe))
			c.Abort()
			return
		}

		if strings.TrimSpace(workflow.WorkflowName) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow name is required"})
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowNameRequired))
			c.Abort()
			return
		}
		if len(workflow.WorkflowName) > 255 { // || !regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString(workflow.WorkflowName)
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowNameInvalid))
			c.Abort()
			return
		}

		if strings.TrimSpace(workflow.DirectoryToSave) == "" {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDirectorySaveRequired))
			c.Abort()
			return
		}
		if len(workflow.DirectoryToSave) > 255 { // || !regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString(workflow.DirectoryToSave)
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDirectoryInvalid))
			c.Abort()
			return
		}

		if workflow.UUID.String() == "" { // && !regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`).MatchString(workflow.UUID)
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UUIDInvalid))
			c.Abort()
			return
		}

		if workflow.CreatedAt != "" && len(workflow.CreatedAt) > 30 {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDateInvalid))
			c.Abort()
			return
		}

		if workflow.UpdatedAt != "" && len(workflow.UpdatedAt) > 30 {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDateInvalid))
			c.Abort()
			return
		}

		c.Set("workflow", workflow)
		c.Next()
	}
}

func ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var currentUser models.SyncUserRequest
		if err := c.ShouldBindJSON(&currentUser); err != nil {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON))
			c.Abort()
			return
		}

		// Validación de `sub`
		if strings.TrimSpace(currentUser.Sub) == "" {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserSubIsRequired))
			c.Abort()
			return
		}
		if len(currentUser.Sub) > 50 { // || !regexp.MustCompile(`^\d+$`).MatchString(currentUser.Sub)
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserSubInvalid))
			c.Abort()
			return
		}

		// Validación de `access_token`
		if currentUser.AccessToken != "" && (len(currentUser.AccessToken) > 1000 || len(currentUser.AccessToken) < 600) { // || !regexp.MustCompile(`^[A-Za-z0-9._-]+$`).MatchString(currentUser.AccessToken))
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserAccessTokenInvalid))
			c.Abort()
			return
		}

		if currentUser.Status != 0 {
			if _, err := models.UserStatusFromUint8(uint8(currentUser.Status)); err != nil {
				c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserInvalidStatus))
				c.Abort()
				return
			}
		}

		if currentUser.RoleID != 0 {
			if currentUser.RoleID < models.RoleAdmin || currentUser.RoleID > models.RoleDeveloper {
				c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserInvalidRole))
				c.Abort()
				return
			}
		}

		c.Set("user", currentUser)
		c.Next()
	}
}
