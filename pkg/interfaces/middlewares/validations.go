package middlewares

import (
	"minireipaz/pkg/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateOnCreateWorkflow() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var workflow models.WorkflowFrontend
		if err := ctx.ShouldBindJSON(&workflow); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if !validateSub(workflow.UserID, ctx) {
			return
		}

		if !validateWorkflowName(workflow.WorkflowName, ctx) {
			return
		}

		if !validateDirectoryToSave(workflow.DirectoryToSave, ctx) {
			return
		}

		// if !validateUUID(workflow.UUID, ctx) {
		// 	return
		// }

		if !validateDates(workflow.CreatedAt, workflow.UpdatedAt, ctx) {
			return
		}

		ctx.Set("workflow", workflow)
		ctx.Next()
	}
}

func ValidateOnGetWorkflow() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: validations
		ctx.Next()
	}
}

func ValidateOnUpdateWorkflow() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var workflow models.Workflow
		if err := ctx.ShouldBindJSON(&workflow); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if !validateSub(workflow.UUID, ctx) {
			return
		}

		if !validateWorkflowName(workflow.Name, ctx) {
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
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
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
		ctx.Next()
	}
}

func ValidateOnCreateCredential() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var currentReq models.RequestCreateCredential
		if err := ctx.ShouldBindJSON(&currentReq); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		ctx.Set(models.CredentialCreateContextKey, currentReq)
		ctx.Next()
	}
}

func ValidateOnExchangeCredential() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var currentReq models.RequestExchangeCredential
		if err := ctx.ShouldBindJSON(&currentReq); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		ctx.Set(models.CredentialExchangeContextKey, currentReq)
		ctx.Next()
	}
}

func ValidateGetGoogleSheet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var currentReq models.RequestGoogleAction
		if err := ctx.ShouldBindBodyWithJSON(&currentReq); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		ctx.Set(models.ActionGoogleKey, currentReq)
		ctx.Next()
	}
}
