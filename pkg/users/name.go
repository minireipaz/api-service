package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUsers(c *gin.Context) {
	c.String(http.StatusOK, "User: %v", c.Param("name"))
}
