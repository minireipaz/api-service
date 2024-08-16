package controllers

import (
	// "fmt"
	// "minireipaz/pkg/domain/models"
	"log"
	"minireipaz/pkg/domain/services"

	"github.com/gin-gonic/gin"
	// "net/http"
	// "github.com/gin-gonic/gin"
)

type DashboardController struct {
	service     *services.DashboardService
	authService *services.AuthService
}

func NewDashboardController(service *services.DashboardService, authServ *services.AuthService) *DashboardController {
	return &DashboardController{
		service:     service,
		authService: authServ,
	}
}

func (d *DashboardController) GetUserDashboardByID(c *gin.Context) {
  id := c.Param("id")
  log.Printf("%v", id)
}
