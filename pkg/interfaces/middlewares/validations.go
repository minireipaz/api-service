package middlewares

import (
	"minireipaz/pkg/domain/models"
	"net/http"

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

		if !validateSub(workflow.Sub, c) {
			return
		}

		if !validateWorkflowName(workflow.WorkflowName, c) {
			return
		}

		if !validateDirectoryToSave(workflow.DirectoryToSave, c) {
			return
		}

		if !validateUUID(workflow.UUID.String(), c) {
			return
		}

		if !validateDates(workflow.CreatedAt, workflow.UpdatedAt, c) {
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

		if !validateSub(currentUser.Sub, c) {
			return
		}

		if !validateAccessToken(currentUser.AccessToken, c) {
			return
		}

		if !validateUserStatus(currentUser.Status, c) {
			return
		}

		if !validateUserRole(currentUser.RoleID, c) {
			return
		}

		c.Set("user", currentUser)
		c.Next()
	}
}

func ValidateUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var currentUser models.SyncUserRequest
		// if err := c.ShouldBindJSON(&currentUser); err != nil {
		// 	c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON))
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
