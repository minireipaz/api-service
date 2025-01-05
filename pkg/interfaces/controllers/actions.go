package controllers

import (
	"fmt"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActionsController struct {
	actionsService repos.ActionsService
	authService    repos.AuthService
}

func NewActionsController(newActionsService repos.ActionsService, newAuthServ repos.AuthService) *ActionsController {
	return &ActionsController{actionsService: newActionsService, authService: newAuthServ}
}

func (a *ActionsController) CreateActionsGoogleSheet(ctx *gin.Context) {
	newAction := ctx.MustGet(models.ActionGoogleKey).(models.RequestGoogleAction)
	actionUserToken, err := a.authService.GetActionUserAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  fmt.Sprintf("Failed to authenticate: %v", err),
			"status": http.StatusInternalServerError,
		})
		return
	}
	sendedBroker, sendedToService, actionID := a.actionsService.CreateActionsGoogleSheet(newAction, actionUserToken)
	if !sendedBroker {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.WorkflowNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}

	if !sendedBroker && !sendedToService { // happens
		ctx.JSON(http.StatusAlreadyReported, gin.H{
			"error":  models.WorkflowNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseGetGoogleSheetByID{
		Status: http.StatusOK,
		Error:  "",
		Data:   *actionID,
	})
}

// func (a *ActionsController) GetGoogleSheetByID(ctx *gin.Context) {
// 	actionID := ctx.Param("idaction")
// 	userID := ctx.Param("iduser")
// 	data, err := a.actionsService.GetGoogleSheetByID(&actionID, &userID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error":  fmt.Sprintf("Failed to authenticate: %v", err),
// 			"status": http.StatusInternalServerError,
// 		})
// 		return
// 	}
// 	if data != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error":  models.WorkflowNameNotGenerate,
// 			"status": http.StatusInternalServerError,
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, models.ResponseGetGoogleSheetByID{
// 		Status: http.StatusOK,
// 		Error:  "",
// 		Data:   *data,
// 	})
// }
