package tests

import (
	"errors"
	"fmt"
	"minireipaz/mocks"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/interfaces/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialController_CreateCredential(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setupMocks         func(*mocks.CredentialService)
		setupContext       func(*gin.Context)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success - Credential Created",
			setupMocks: func(mockCredService *mocks.CredentialService) {
				mockCredService.On(
					"CreateCredential",
					mock.MatchedBy(func(cred *models.RequestCreateCredential) bool {
						return cred.Sub == "test-sub" && cred.ID == "test-id" && cred.Type == "googlesheets"
					}),
				).Return(&models.RequestExchangeCredential{
					Data: models.DataCredential{
						RedirectURL: "https://oauth.google.com/auth",
					},
				}, nil)
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.CredentialCreateContextKey, models.RequestCreateCredential{
					Sub:  "test-sub",
					ID:   "test-id",
					Type: "googlesheets",
				})
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"auth_url":"https://oauth.google.com/auth","error":"","status":200}`,
		},
		{
			name: "Error - Invalid Credential Type",
			setupMocks: func(mockCredService *mocks.CredentialService) {
				mockCredService.On(
					"CreateCredential",
					mock.AnythingOfType("*models.RequestCreateCredential"),
				).Return(nil, errors.New("invalid credential type"))
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.CredentialCreateContextKey, models.RequestCreateCredential{
					Sub:  "test-sub",
					ID:   "test-id",
					Type: "unsupported",
				})
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   fmt.Sprintf(`{"error": "%s","status": 500}`, models.CredNameNotGenerate),
		},
		{
			name: "Error - Service Internal Failure",
			setupMocks: func(mockCredService *mocks.CredentialService) {
				mockCredService.On(
					"CreateCredential",
					mock.AnythingOfType("*models.RequestCreateCredential"),
				).Return(nil, errors.New("internal service failure"))
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.CredentialCreateContextKey, models.RequestCreateCredential{
					Sub:  "test-sub",
					ID:   "test-id",
					Type: "googlesheets",
				})
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   fmt.Sprintf(`{"error": "%s","status": 500}`, models.CredNameNotGenerate),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCredService := mocks.NewCredentialService(t)
			tt.setupMocks(mockCredService)
			controller := controllers.NewCredentialController(mockCredService, nil, nil)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			tt.setupContext(ctx)

			controller.CreateCredential(ctx)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponse, w.Body.String())

			mockCredService.AssertExpectations(t)
		})
	}
}
