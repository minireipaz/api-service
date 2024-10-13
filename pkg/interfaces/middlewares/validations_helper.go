package middlewares

import (
	"minireipaz/pkg/domain/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateSub(sub string, ctx *gin.Context) bool {
	if strings.TrimSpace(sub) == "" || len(sub) > 50 {
		ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserSubMustBe))
		ctx.Abort()
		return false
	}
	return true
}

func validateWorkflowName(name string, ctx *gin.Context) bool {
	if strings.TrimSpace(name) == "" || len(name) > 255 {
		ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowNameInvalid))
		ctx.Abort()
		return false
	}
	return true
}

func validateDirectoryToSave(directory string, ctx *gin.Context) bool {
	if strings.TrimSpace(directory) == "" || len(directory) > 255 {
		ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDirectoryInvalid))
		ctx.Abort()
		return false
	}
	return true
}

func validateUUID(uuid string, ctx *gin.Context) bool {
	if strings.TrimSpace(uuid) == " " { // init validation, right now cannot
		ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UUIDInvalid))
		ctx.Abort()
		return false
	}
	return true
}

func validateDates(createdAt, updatedAt string, ctx *gin.Context) bool {
	if (createdAt != "" && len(createdAt) > 30) || (updatedAt != "" && len(updatedAt) > 30) {
		ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDateInvalid))
		ctx.Abort()
		return false
	}
	return true
}

func validateAccessToken(token string, ctx *gin.Context) bool {
	if token != "" && (len(token) > 1000 || len(token) < 100) {
		ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserAccessTokenInvalid))
		ctx.Abort()
		return false
	}
	return true
}

func validateUserStatus(status models.UserStatus, ctx *gin.Context) bool {
	if status != 0 {
		if _, err := models.UserStatusFromUint8(uint8(status)); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserInvalidStatus))
			ctx.Abort()
			return false
		}
	}
	return true
}

func validateUserRole(roleID models.UserRoleID, ctx *gin.Context) bool {
	if roleID != 0 {
		if roleID < models.RoleAdmin || roleID > models.RoleDeveloper {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserInvalidRole))
			ctx.Abort()
			return false
		}
	}
	return true
}
