package controllers

import (
	"fmt"
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service     repos.DashboardService
	authService repos.AuthService
}

func NewDashboardController(service repos.DashboardService, authServ repos.AuthService) *DashboardController {
	return &DashboardController{
		service:     service,
		authService: authServ,
	}
}

func (d *DashboardController) GetUserDashboardByID(ctx *gin.Context) {
	id := ctx.Param("iduser")
	dashboardInfo, err := d.service.QueryDashboardInfo(id)
	resp := models.ResponseInfoDashboard{
		Data:   dashboardInfo.Data,
		Status: http.StatusOK,
		Error:  "",
	}
	if err != nil {
		log.Printf("ERROR | %v", err)
		resp.Status = http.StatusInternalServerError
		resp.Error = fmt.Sprintf("ERROR | %v", err)
	}
	ctx.JSON(resp.Status, resp)
}
