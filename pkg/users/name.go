package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUsers(ctx *gin.Context) {
	ctx.String(http.StatusOK, "User: %v", ctx.Param("name"))
}
