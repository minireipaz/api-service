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
	workflowFrontend := ctx.MustGet("workflow").(models.WorkflowFrontend)
	created, exist := c.workflowService.CreateWorkflow(&workflowFrontend)
	if !created && !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    models.WorkflowNameNotGenerate,
			"workflow": workflowFrontend,
			"status":   http.StatusInternalServerError,
		})
		return
	}

	if exist {
		ctx.JSON(http.StatusAlreadyReported, gin.H{
			"error":    models.WorkflowNameExist,
			"workflow": workflowFrontend,
			"status":   http.StatusAlreadyReported,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":    "",
		"workflow": workflowFrontend,
		"status":   http.StatusCreated,
	})
}
