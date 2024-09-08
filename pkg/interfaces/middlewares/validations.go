package middlewares

import (
	"minireipaz/pkg/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateWorkflow() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var workflow models.WorkflowFrontend
		if err := ctx.ShouldBindJSON(&workflow); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON))
			ctx.Abort()
			return
		}

		if !validateSub(workflow.Sub, ctx) {
			return
		}

		if !validateWorkflowName(workflow.WorkflowName, ctx) {
			return
		}

		if !validateDirectoryToSave(workflow.DirectoryToSave, ctx) {
			return
		}

		if !validateUUID(workflow.UUID, ctx) {
			return
		}

		if !validateDates(workflow.CreatedAt, workflow.UpdatedAt, ctx) {
			return
		}

		ctx.Set("workflow", workflow)
		ctx.Next()
	}
}

func ValidateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var currentUser models.SyncUserRequest
		if err := ctx.ShouldBindJSON(&currentUser); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON))
			ctx.Abort()
			return
		}

		if !validateSub(currentUser.Sub, ctx) {
			return
		}

		if !validateAccessToken(currentUser.AccessToken, ctx) {
			return
		}

		if !validateUserStatus(currentUser.Status, ctx) {
			return
		}

		if !validateUserRole(currentUser.RoleID, ctx) {
			return
		}

		ctx.Set("user", currentUser)
		ctx.Next()
	}
}

func ValidateUserAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var currentUser models.SyncUserRequest
		// if err := c.ShouldBindJSON(&currentUser); err != nil {
		// 	c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON))
		// 	c.Abort()
		// 	return
		// }
		ctx.Next()
	}
}
