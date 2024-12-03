package controllers

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActionsController struct {
	actionsService repos.ActionsService
}

func NewActionsController(newActionsService repos.ActionsService) *ActionsController {
	return &ActionsController{actionsService: newActionsService}
}

func (a *ActionsController) GetGoogleSheetByID(ctx *gin.Context) {
	newAction := ctx.MustGet(models.ActionGoogleKey).(models.RequestGoogleAction)
	created, exist, actionsData := a.actionsService.GetGoogleSheetByID(newAction)
	if !created && !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.WorkflowNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}

	if exist {
		ctx.JSON(http.StatusAlreadyReported, gin.H{
			"error":  models.WorkflowNameExist,
			"status": http.StatusAlreadyReported,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseGetGoogleSheetByID{
		Status: http.StatusOK,
		Error:  "",
		Action: *actionsData,
	})
}
