package middlewares

import (
	"minireipaz/pkg/domain/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateSub(sub string, c *gin.Context) bool {
	if strings.TrimSpace(sub) == "" || len(sub) > 50 {
		c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserSubMustBe))
		c.Abort()
		return false
	}
	return true
}

func validateWorkflowName(name string, c *gin.Context) bool {
	if strings.TrimSpace(name) == "" || len(name) > 255 {
		c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowNameInvalid))
		c.Abort()
		return false
	}
	return true
}

func validateDirectoryToSave(directory string, c *gin.Context) bool {
	if strings.TrimSpace(directory) == "" || len(directory) > 255 {
		c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDirectoryInvalid))
		c.Abort()
		return false
	}
	return true
}

func validateUUID(uuid string, c *gin.Context) bool {
	if uuid == "" {
		c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UUIDInvalid))
		c.Abort()
		return false
	}
	return true
}

func validateDates(createdAt, updatedAt string, c *gin.Context) bool {
	if (createdAt != "" && len(createdAt) > 30) || (updatedAt != "" && len(updatedAt) > 30) {
		c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.WorkflowDateInvalid))
		c.Abort()
		return false
	}
	return true
}

func validateAccessToken(token string, c *gin.Context) bool {
	if token != "" && (len(token) > 1000 || len(token) < 100) {
		c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserAccessTokenInvalid))
		c.Abort()
		return false
	}
	return true
}

func validateUserStatus(status models.UserStatus, c *gin.Context) bool {
	if status != 0 {
		if _, err := models.UserStatusFromUint8(uint8(status)); err != nil {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserInvalidStatus))
			c.Abort()
			return false
		}
	}
	return true
}

func validateUserRole(roleID models.UserRoleID, c *gin.Context) bool {
	if roleID != 0 {
		if roleID < models.RoleAdmin || roleID > models.RoleDeveloper {
			c.JSON(http.StatusBadRequest, NewInvalidRequestError(models.UserInvalidRole))
			c.Abort()
			return false
		}
	}
	return true
}
