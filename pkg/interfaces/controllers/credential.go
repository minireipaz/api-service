package controllers

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CredentialController struct {
	credentialService repos.CredentialService
	authService       repos.AuthService
}

func NewCredentialController(credServ repos.CredentialService, authService repos.AuthService) *CredentialController {
	return &CredentialController{credentialService: credServ, authService: authService}
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
	token, tokenRefresh, _, stateInfo, err := c.credentialService.ExchangeGoogleCredential(&currentCredential)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  models.CredNameNotGenerate,
			"status": http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":        "",
		"status":       http.StatusOK,
		"token":        token,
		"tokenrefresh": tokenRefresh,
		"id":           stateInfo.ID,
	})
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
