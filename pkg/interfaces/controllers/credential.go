package controllers

import (
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CredentialController struct {
	credentialService repos.CredentialService
	authService       repos.AuthService
	workflowService   repos.WorkflowService
}

func NewCredentialController(credServ repos.CredentialService, authService repos.AuthService, workflowServ repos.WorkflowService) *CredentialController {
	return &CredentialController{credentialService: credServ, authService: authService, workflowService: workflowServ}
}

func (c *CredentialController) CreateCredential(ctx *gin.Context) {
	credFrontend := ctx.MustGet(models.CredentialCreateContextKey).(models.RequestCreateCredential)
	currentCredential, err := c.credentialService.CreateCredential(&credFrontend)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.CredNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"auth_url": currentCredential.Data.RedirectURL,
		"error":    "",
		"status":   http.StatusOK,
	})
}

func (c *CredentialController) ExchangeGoogleCode(ctx *gin.Context) {
	currentCredential := ctx.MustGet(models.CredentialExchangeContextKey).(models.RequestExchangeCredential)
	// this token and refresh expire in 1hr
	// stateinfo all returned because dont know what values are necessary in this controller
	// will lock workflow
	token, tokenRefresh, _, stateInfo, err := c.credentialService.ExchangeGoogleCredential(&currentCredential)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.CredNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}
	done := make(chan bool)
	// using goroutine for quick-response
	// update workflow only if token generated????
	// TODO: better consistency atomicity
	if token != nil {
		go func() {
			workflow, exist := c.workflowService.GetWorkflow(&currentCredential.Sub, &currentCredential.WorkflowID)
			// its necessary to send Notfound status when workflow not exist??
			if !exist {
				log.Printf("ERROR | workflow not found for sub: %s and workflowid: %s", currentCredential.Sub, currentCredential.WorkflowID)
				done <- true // goroutine ended
			}
			workflow = c.credentialService.TransformWorkflow(&currentCredential, workflow)
			// block workflow key with retries
			// max retries 10
			updated, _ := c.workflowService.UpdateWorkflow(workflow)
			if !updated {
				log.Printf("ERROR | failed to update workflow: %s", workflow.UUID)
				// TODO: dead letter
			}
			done <- true // goroutine ended
		}()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":        "",
		"status":       http.StatusOK,
		"token":        token,
		"tokenrefresh": tokenRefresh,
		"id":           stateInfo.ID,
	})

	select {
	case <-done:
	case <-time.After(models.MaxSecondsGoRoutine):
		log.Println("WARN | Needs more time than 10 seconds")
	}
}

func (c *CredentialController) GetAllCredentials(ctx *gin.Context) {
	userID := ctx.Param("iduser")
	credentials, exist := c.credentialService.GetAllCredentials(&userID)

	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":  models.UUIDInvalid,
			"status": http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, credentials)
}

func (c *CredentialController) CreateTokenCredential(ctx *gin.Context) {
	credFrontend := ctx.MustGet(models.CredentialCreateContextKey).(models.RequestCreateCredential)
	savedCredential, transformedCredentialID, err := c.credentialService.CreateTokenCredential(&credFrontend)
	if err != nil || !savedCredential {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.CredNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error":  "",
		"status": http.StatusOK,
		"id":     *transformedCredentialID,
	})
}
