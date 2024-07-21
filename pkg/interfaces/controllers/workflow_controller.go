package controllers

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkflowController struct {
	workflowService *services.WorkflowService
	authService     *services.AuthService
}

func NewWorkflowController(newWorkflowService *services.WorkflowService, newAuthService *services.AuthService) *WorkflowController {
	return &WorkflowController{workflowService: newWorkflowService, authService: newAuthService}
}

func (c *WorkflowController) CreateWorkflow(ctx *gin.Context) {
	workflow := ctx.MustGet("workflow").(models.Workflow)
	created, exist := c.workflowService.CreateWorkflow(&workflow)
	if !created && !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    models.WorkflowNameNotGenerate,
			"workflow": workflow,
			"status":   http.StatusInternalServerError,
		})
		return
	}

	if exist {
		ctx.JSON(http.StatusAlreadyReported, gin.H{
			"error":    models.WorkflowNameExist,
			"workflow": workflow,
			"status":   http.StatusAlreadyReported,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":    "",
		"workflow": workflow,
		"status":   http.StatusCreated,
	})
}
