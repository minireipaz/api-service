package controllers

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkflowController struct {
	authService       repos.AuthService
	workflowService   repos.WorkflowService
	credentialService repos.CredentialService
}

func NewWorkflowController(newWorkflowService repos.WorkflowService, newCredentialService repos.CredentialService, newAuthService repos.AuthService) *WorkflowController {
	return &WorkflowController{
		workflowService:   newWorkflowService,
		credentialService: newCredentialService,
		authService:       newAuthService,
	}
}

func (c *WorkflowController) CreateWorkflow(ctx *gin.Context) {
	workflowFrontend := ctx.MustGet("workflow").(models.WorkflowFrontend)
	created, exist, workflow := c.workflowService.CreateWorkflow(&workflowFrontend)
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
		"workflow": workflow,
		"status":   http.StatusCreated,
	})
}

func (c *WorkflowController) GetWorkflow(ctx *gin.Context) {
	userID := ctx.Param("iduser")
	workflowID := ctx.Param("idworkflow")

	newWorkflow, exist := c.workflowService.GetWorkflow(&userID, &workflowID)

	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":  models.UUIDInvalid,
			"status": http.StatusNotFound,
		})
		return
	}

	credentials, _ := c.credentialService.GetAllCredentials(&userID)
	if credentials.Status != 200 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":  models.UUIDInvalid,
			"status": http.StatusNotFound,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error":       "",
		"status":      http.StatusOK,
		"workflow":    newWorkflow,
		"credentials": credentials.Credentials,
	})
}

func (c *WorkflowController) GetAllWorkflows(ctx *gin.Context) {
	userID := ctx.Param("iduser")
	allWorkflows, err := c.workflowService.GetAllWorkflows(&userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":  models.UUIDInvalid,
			"status": http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":    "",
		"status":   http.StatusOK,
		"workflow": allWorkflows,
	})
}

func (c *WorkflowController) UpdateWorkflow(ctx *gin.Context) {
	workflowFrontend := ctx.MustGet("workflow").(models.Workflow)
	updated, exist := c.workflowService.UpdateWorkflow(&workflowFrontend)

	if !exist {
		ctx.JSON(http.StatusAlreadyReported, gin.H{
			"error":  models.WorkflowNameExist,
			"status": http.StatusNotFound,
		})
		return
	}

	if !updated {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.WorkflowNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":  "",
		"status": http.StatusOK,
	})
}
